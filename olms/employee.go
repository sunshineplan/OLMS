package olms

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type employee struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Realname   string `json:"realname"`
	DeptID     int    `json:"deptid"`
	DeptName   string `json:"deptname"`
	Role       bool   `json:"role"`
	Permission string `json:"permission"`
}

func getUser(db *sql.DB, id interface{}) (user employee, err error) {
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
	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	deptID := c.PostForm("dept")
	if deptID != "" && !checkPermission(db, c, &idOptions{Departments: []string{fmt.Sprintf("%v", deptID)}}) {
		c.String(403, "")
		return
	}
	username := strings.TrimSpace(c.PostForm("username"))
	realname := strings.TrimSpace(c.PostForm("realname"))
	if realname == "" {
		realname = username
	}
	localize := localize(c)
	var exist, message string
	var code int
	if username == "" {
		message = "Username is required."
	} else if err := db.QueryRow(
		"SELECT id FROM user WHERE username = ?", strings.ToLower(username)).Scan(&exist); err == nil {
		message = fmt.Sprintf(localize["UsernameExist"], username)
		code = 1
	} else if deptID == "" {
		message = localize["DepartmentRequired"]
	} else {
		if checkSuper(c) {
			role := c.PostForm("role")
			res, err := db.Exec("INSERT INTO user (username, realname, dept_id, role) VALUES (?, ?, ?, ?)",
				strings.ToLower(username), realname, deptID, role)
			if err != nil {
				log.Println("Failed to add user:", err)
				c.String(500, "")
				return
			}
			id, err := res.LastInsertId()
			if err != nil {
				log.Println("Failed to get last insert id:", err)
				c.String(500, "")
				return
			}
			if role == "1" {
				for _, permission := range c.PostFormArray("permission") {
					if _, err = db.Exec(
						"INSERT INTO permission (dept_id, user_id) VALUES (?, ?)", permission, id); err != nil {
						log.Println("Failed to add user permission:", err)
						c.String(500, "")
						return
					}
				}
			}
			c.JSON(200, gin.H{"status": 1})
			return
		}
		if _, err = db.Exec(
			"INSERT INTO user (username, realname, dept_id) VALUES (?, ?, ?)", username, realname, deptID); err != nil {
			log.Println("Failed to add user:", err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		return
	}
	c.JSON(200, gin.H{"status": 0, "message": message, "error": code})
}

func editEmployee(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	deptID := c.PostForm("dept")
	if deptID != "" && !checkPermission(db, c, &idOptions{Departments: []string{fmt.Sprintf("%v", deptID)}}) {
		c.String(403, "")
		return
	}
	id := c.Param("id")
	username := strings.TrimSpace(c.PostForm("username"))
	realname := strings.TrimSpace(c.PostForm("realname"))
	if realname == "" {
		realname = username
	}
	role := c.PostForm("role")

	localize := localize(c)
	var exist, message string
	if username == "" {
		message = "Username is required."
	} else if err := db.QueryRow("SELECT id FROM user WHERE username = ? AND id != ?",
		strings.ToLower(username), id).Scan(&exist); err == nil {
		message = fmt.Sprintf(localize["UsernameExist"], username)
	} else if deptID == "" {
		message = localize["DepartmentRequired"]
	} else {
		if password := c.PostForm("password"); password == "" {
			if _, err = db.Exec(
				"UPDATE user SET username = ?, realname = ?, dept_id = ?, role = ? WHERE id = ?",
				strings.ToLower(username), realname, deptID, role, id); err != nil {
				log.Println("Failed to edit user:", err)
				c.String(500, "")
				return
			}
		} else {
			newPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				log.Print(err)
				c.String(500, "")
				return
			}
			if _, err = db.Exec(
				"UPDATE user SET username = ?, realname = ?, password = ?, dept_id = ?, role = ? WHERE id = ?",
				strings.ToLower(username), realname, string(newPassword), deptID, role, id); err != nil {
				log.Println("Failed to edit user:", err)
				c.String(500, "")
				return
			}
		}
		if _, err = db.Exec("DELETE FROM permission WHERE user_id = ?", id); err != nil {
			log.Println("Failed to clear user permission:", err)
			c.String(500, "")
			return
		}
		if role == "1" {
			for _, permission := range c.PostFormArray("permission") {
				if _, err = db.Exec(
					"INSERT INTO permission (dept_id, user_id) VALUES (?, ?)", permission, id); err != nil {
					log.Println("Failed to add user permission:", err)
					c.String(500, "")
					return
				}
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
	if !checkPermission(db, c, &idOptions{User: id}) {
		c.String(403, "")
		return
	}
	if _, err := db.Exec("DELETE FROM user WHERE id = ?", id); err != nil {
		log.Println("Failed to delete employee:", err)
		c.String(500, "")
		return
	}
	c.JSON(200, gin.H{"status": 1})
}
