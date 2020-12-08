package olms

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type employee struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Realname   string `json:"realname"`
	Password   string
	DeptID     int    `json:"deptid"`
	DeptName   string `json:"deptname"`
	Role       bool   `json:"role"`
	Permission string `json:"permission"`
}

func getUser(db *sql.DB, id interface{}) (user employee, err error) {
	if id == "0" {
		user = employee{ID: 0, Realname: "root", Role: true}
		return
	}
	var permission []byte
	if err = db.QueryRow("SELECT * FROM employee WHERE id = ?", id).Scan(
		&user.ID, &user.Username, &user.Realname, &user.DeptID, &user.DeptName, &user.Role, &permission); err != nil {
		log.Println("Failed to get user:", err)
		return
	}
	user.Permission = string(permission)
	return
}

func getEmployees(db *sql.DB, ids []string, super bool) (employees []employee, err error) {
	var rows *sql.Rows
	if super {
		rows, err = db.Query("SELECT * FROM employee")
		if err != nil {
			log.Println("Failed to get employees:", err)
			return
		}
		for rows.Next() {
			var e employee
			var permission []byte
			if err = rows.Scan(&e.ID, &e.Username, &e.Realname, &e.DeptID, &e.DeptName, &e.Role, &permission); err != nil {
				log.Println("Failed to scan department:", err)
				return
			}
			e.Permission = string(permission)
			employees = append(employees, e)
		}
	} else {
		rows, err = db.Query("SELECT * FROM employee WHERE id IN (" + strings.Join(ids, ", ") + ")")
		if err != nil {
			log.Println("Failed to get employees:", err)
			return
		}
		for rows.Next() {
			var e employee
			if err = rows.Scan(&e.ID, &e.Username, &e.Realname, &e.DeptID, &e.DeptName); err != nil {
				log.Println("Failed to scan department:", err)
				return
			}
			employees = append(employees, e)
		}
	}
	return
}

func addEmployee(c *gin.Context) {
	var employee employee
	if err := c.BindJSON(&employee); err != nil {
		c.String(400, "")
		return
	}

	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	if employee.DeptID != 0 &&
		!checkPermission(db, c, &idOptions{Departments: []string{strconv.Itoa(employee.DeptID)}}) {
		c.String(403, "")
		return
	}
	if employee.Realname == "" {
		employee.Realname = employee.Username
	}

	var exist, message string
	var code int
	if employee.Username == "" {
		message = "UsernameRequired"
	} else if err := db.QueryRow(
		"SELECT id FROM user WHERE username = ?", strings.ToLower(employee.Username)).Scan(&exist); err == nil {
		message = "UsernameExist"
		code = 1
	} else if employee.DeptID == 0 {
		message = "DepartmentRequired"
	} else {
		if sessions.Default(c).Get("userID") == "0" {
			res, err := db.Exec("INSERT INTO user (username, realname, dept_id) VALUES (?, ?, ?, ?)",
				strings.ToLower(employee.Username), employee.Realname, employee.DeptID)
			if err != nil {
				log.Println("Failed to add employee:", err)
				c.String(500, "")
				return
			}
			id, err := res.LastInsertId()
			if err != nil {
				log.Println("Failed to get last insert id:", err)
				c.String(500, "")
				return
			}
			if employee.Role {
				if employee.Permission != "" {
					if _, err := db.Exec("UPDATE user SET role = 1 WHERE id = ?", id); err != nil {
						log.Println("Failed to save employee role:", err)
						c.String(500, "")
						return
					}
					for _, permission := range strings.Split(employee.Permission, ",") {
						if _, err := db.Exec(
							"INSERT INTO permission (dept_id, user_id) VALUES (?, ?)", permission, id); err != nil {
							log.Println("Failed to add employee permission:", err)
							c.String(500, "")
							return
						}
					}
				} else {
					c.JSON(200, gin.H{"status": 0, "message": "EmptyPermission"})
					return
				}
			}
			c.JSON(200, gin.H{"status": 1})
			return
		}
		if _, err := db.Exec(
			"INSERT INTO user (username, realname, dept_id) VALUES (?, ?, ?)",
			employee.Username, employee.Realname, employee.DeptID); err != nil {
			log.Println("Failed to add employee:", err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		return
	}
	c.JSON(200, gin.H{"status": 0, "message": message, "error": code})
}

func editEmployee(c *gin.Context) {
	var employee employee
	if err := c.BindJSON(&employee); err != nil {
		c.String(400, "")
		return
	}

	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	if employee.DeptID != 0 &&
		!checkPermission(db, c, &idOptions{Departments: []string{strconv.Itoa(employee.DeptID)}}) {
		c.String(403, "")
		return
	}
	if employee.Realname == "" {
		employee.Realname = employee.Username
	}

	var exist, message string
	if employee.Username == "" {
		message = "UsernameRequired"
	} else if err := db.QueryRow("SELECT id FROM user WHERE username = ? AND id != ?",
		strings.ToLower(employee.Username), employee.ID).Scan(&exist); err == nil {
		message = "UsernameExist"
	} else if employee.DeptID == 0 {
		message = "DepartmentRequired"
	} else {
		if employee.Password == "" {
			if _, err := db.Exec(
				"UPDATE user SET username = ?, realname = ?, dept_id = ? WHERE id = ?",
				strings.ToLower(employee.Username), employee.Realname, employee.DeptID, employee.ID); err != nil {
				log.Println("Failed to edit employee:", err)
				c.String(500, "")
				return
			}
		} else {
			newPassword, err := bcrypt.GenerateFromPassword([]byte(employee.Password), bcrypt.MinCost)
			if err != nil {
				log.Print(err)
				c.String(500, "")
				return
			}
			if _, err := db.Exec(
				"UPDATE user SET username = ?, realname = ?, password = ?, dept_id = ? WHERE id = ?",
				strings.ToLower(employee.Username), employee.Realname, string(newPassword), employee.DeptID, employee.ID,
			); err != nil {
				log.Println("Failed to edit employee:", err)
				c.String(500, "")
				return
			}
		}
		if _, err := db.Exec("DELETE FROM permission WHERE user_id = ?", employee.ID); err != nil {
			log.Println("Failed to clear user permission:", err)
			c.String(500, "")
			return
		}
		if employee.Role {
			if employee.Permission != "" {
				if _, err := db.Exec("UPDATE user SET role = 1 WHERE id = ?", employee.ID); err != nil {
					log.Println("Failed to save employee role:", err)
					c.String(500, "")
					return
				}
				for _, permission := range strings.Split(employee.Permission, ",") {
					if _, err := db.Exec(
						"INSERT INTO permission (dept_id, user_id) VALUES (?, ?)", permission, employee.ID); err != nil {
						log.Println("Failed to save employee permission:", err)
						c.String(500, "")
						return
					}
				}
			} else {
				c.JSON(200, gin.H{"status": 0, "message": "EmptyPermission"})
				return
			}
		}
		c.JSON(200, gin.H{"status": 1})
		return
	}
	c.JSON(200, gin.H{"status": 0, "message": message})
}

func deleteEmployee(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	id := c.Param("id")
	if _, err := db.Exec("DELETE FROM user WHERE id = ?", id); err != nil {
		log.Println("Failed to delete employee:", err)
		c.String(500, "")
		return
	}
	c.JSON(200, gin.H{"status": 1})
}
