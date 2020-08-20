package olms

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type empl struct {
	ID         int
	Username   string
	Realname   string
	DeptID     int
	DeptName   string
	Role       bool
	Permission string
}

func getEmpls(id *idOptions, options *searchOptions) (empls []empl, total int, err error) {
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	stmt := "SELECT %s FROM employee WHERE"

	var args []interface{}
	var orderBy, limit string
	bc := make(chan bool, 1)
	if id.UserID != nil {
		stmt += " id = ?"
		args = append(args, id.UserID)
		bc <- true
	} else {
		marks := make([]string, len(id.DeptIDs))
		for i := range marks {
			marks[i] = "?"
		}
		stmt += " dept_id IN (" + strings.Join(marks, ", ") + ")"
		for _, i := range id.DeptIDs {
			args = append(args, i)
		}

		if options != nil {
			if r, ok := options.Role.(float64); ok {
				stmt += " AND role = ?"
				args = append(args, r)
			}
			if p, ok := options.Page.(float64); ok {
				limit = fmt.Sprintf(" LIMIT ?, ?")
				args = append(args, int(p-1)*perPage, perPage)
			}
			if options.Sort != nil {
				orderBy = fmt.Sprintf(" ORDER BY %v %v", options.Sort, options.Order)
			} else {
				orderBy = " ORDER BY dept_name, realname"
			}
		}
		go func() {
			if err := db.QueryRow(fmt.Sprintf(stmt, "count(*)"), args...).Scan(&total); err != nil {
				log.Printf("Failed to get total records: %v", err)
				bc <- false
			}
			bc <- true
		}()
	}
	rows, err := db.Query(fmt.Sprintf(stmt+orderBy+limit, "id, username, realname, dept_id, dept_name, role, permission"), args...)
	if err != nil {
		log.Printf("Failed to get employees: %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var empl empl
		if err = rows.Scan(&empl.ID, &empl.Username, &empl.Realname, &empl.DeptID, &empl.DeptName, &empl.Role, &empl.Permission); err != nil {
			log.Printf("Failed to scan employee: %v", err)
			return
		}
		empls = append(empls, empl)
	}
	if v := <-bc; !v {
		err = fmt.Errorf("Failed to get total records")
	}
	return
}

func doAddEmpl(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	deptID := c.PostForm("dept")
	if deptID != "" && !checkPermission(c, deptID) {
		c.String(403, "")
		return
	}
	username := strings.TrimSpace(c.PostForm("username"))
	realname := strings.TrimSpace(c.PostForm("realname"))
	if realname == "" {
		realname = username
	}
	var exist, message string
	var code int
	if username == "" {
		message = "Username is required."
	} else if err := db.QueryRow("SELECT id FROM user WHERE username = ?", strings.ToLower(username)).Scan(&exist); err == nil {
		message = fmt.Sprintf("Username %s is already existed.", username)
		code = 1
	} else if deptID == "" {
		message = "Department is required."
	} else {
		if checkSuper(c) {
			role := c.PostForm("role")
			var permission []string
			if role == "1" {
				permission = c.PostFormArray("permission")
			}
			if _, err = db.Exec("INSERT INTO user (username, realname, dept_id, role, permission) VALUES (?, ?, ?, ?, ?)",
				strings.ToLower(username), realname, deptID, role, strings.Join(permission, ",")); err != nil {
				log.Printf("Failed to add user: %v", err)
				c.String(500, "")
				return
			}
			c.JSON(200, gin.H{"status": 1})
			return
		}
		if _, err = db.Exec("INSERT INTO user (username, realname, dept_id) VALUES (?, ?, ?)", username, realname, deptID); err != nil {
			log.Printf("Failed to add user: %v", err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		return
	}
	c.JSON(200, gin.H{"status": 0, "message": message, "error": code})
}

func doEditEmpl(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	deptID := c.PostForm("dept")
	if deptID != "" && !checkPermission(c, deptID) {
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
	var permission []string
	if role == "1" {
		permission = c.PostFormArray("permission")
	}
	var exist, message string
	if username == "" {
		message = "Username is required."
	} else if err := db.QueryRow("SELECT id FROM user WHERE username = ? AND id != ?",
		strings.ToLower(username), id).Scan(&exist); err == nil {
		message = fmt.Sprintf("Username %s is already existed.", username)
	} else if deptID == "" {
		message = "Department is required."
	} else {
		if password := c.PostForm("password"); password == "" {
			if _, err = db.Exec("UPDATE user SET username = ?, realname = ?, dept_id = ?, role = ?, permission = ? WHERE id = ?",
				strings.ToLower(username), realname, deptID, role, strings.Join(permission, ","), id); err != nil {
				log.Printf("Failed to edit user: %v", err)
				c.String(500, "")
				return
			}
		} else {
			newPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				log.Println(err)
				c.String(500, "")
				return
			}
			if _, err = db.Exec(
				"UPDATE user SET username = ?, realname = ?, password = ?, dept_id = ?, role = ?, permission = ? WHERE id = ?",
				strings.ToLower(username), realname, string(newPassword), deptID, role, strings.Join(permission, ","), id); err != nil {
				log.Printf("Failed to edit user: %v", err)
				c.String(500, "")
				return
			}
		}
		c.JSON(200, gin.H{"status": 1})
		return
	}
	c.JSON(200, gin.H{"status": 0, "message": message})
}

func doDeleteEmpl(c *gin.Context) {
	id := c.Param("id")
	empls, _, err := getEmpls(&idOptions{UserID: id}, nil)
	if err != nil {
		log.Printf("Failed to get empl: %v", err)
		c.String(400, "")
		return
	}
	if !checkPermission(c, strconv.Itoa(empls[0].DeptID)) {
		c.String(403, "")
		return
	}
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.String(503, "")
		return
	}
	defer db.Close()
	if _, err := db.Exec("DELETE FROM user WHERE id = ?", id); err != nil {
		log.Printf("Failed to delete employee: %v", err)
		c.String(500, "")
		return
	}
	c.JSON(200, gin.H{"status": 1})
}
