package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type empl struct {
	ID         int
	Realname   string
	DeptID     int
	DeptName   string
	Admin      bool
	Permission string
}

func getEmpl(id interface{}) (empl, error) {
	var empl empl
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return empl, err
	}
	defer db.Close()
	if err := db.QueryRow("SELECT id, realname, dept_id, dept_name, admin, permission FROM employee WHERE id = ?",
		id).Scan(&empl.ID, &empl.Realname, &empl.DeptID, &empl.DeptName, &empl.Admin, &empl.Permission); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Printf("No results: %v", err)
		} else {
			log.Printf("Failed to query employee: %v", err)
		}
		return empl, err
	}
	return empl, nil
}

func checkPermission(id string, c *gin.Context) bool {
	session := sessions.Default(c)
	user, err := getEmpl(session.Get("userID"))
	if err != nil {
		return false
	}
	for _, i := range strings.Split(user.Permission, ",") {
		if id == i {
			return true
		}
	}
	return false
}

func getDeptEmpls(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	deptID := c.Query("dept")
	if !checkPermission(deptID, c) {
		c.String(403, "")
		return
	}
	var empls []empl
	rows, err := db.Query("SELECT id, realname FROM user where dept_id = ? ORDER BY realname", deptID)
	if err != nil {
		log.Printf("Failed to get employees: %v", err)
		c.String(500, "")
		return
	}
	defer rows.Close()
	for rows.Next() {
		var empl empl
		if err := rows.Scan(&empl.ID, &empl.Realname); err != nil {
			log.Printf("Failed to scan employee: %v", err)
			c.String(500, "")
			return
		}
		empls = append(empls, empl)
	}
	c.JSON(200, empls)
}

func getEmpls(c *gin.Context) {
	// do something ...
}

func addEmpl(c *gin.Context) {
	c.HTML(200, "addEmpl.html", nil)
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
	if deptID != "" && !checkPermission(deptID, c) {
		c.String(403, "")
		return
	}
	var message string
	username := strings.TrimSpace(c.PostForm("username"))
	realname := strings.TrimSpace(c.PostForm("realname"))
	if realname == "" {
		realname = username
	}
	var exist string
	if username == "" {
		message = "Username is required."
	} else if err := db.QueryRow("SELECT id FROM user WHERE username = ?", username).Scan(&exist); err == nil {
		message = fmt.Sprintf("Username %s is already existed.", username)
	} else if deptID == "" {
		message = "Department is required."
	} else {
		if _, err = db.Exec("INSERT INTO user (username, realname, dept_id)' VALUES (?, ?, ?, ?, ?)", username, realname, deptID); err != nil {
			log.Printf("Failed to add user: %v", err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		return
	}
	c.JSON(200, gin.H{"status": 0, "message": message})
}

func editEmpl(c *gin.Context) {
	id := c.Param("id")
	empl, err := getEmpl(id)
	if err != nil {
		log.Printf("Failed to get empl: %v", err)
		c.String(400, "")
		return
	}
	if !checkPermission(strconv.Itoa(empl.DeptID), c) {
		c.String(403, "")
		return
	}
	c.HTML(200, "addDept.html", gin.H{"empl": empl})
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
	if deptID != "" && !checkPermission(deptID, c) {
		c.String(403, "")
		return
	}
	var message string
	id := c.Param("id")
	username := strings.TrimSpace(c.PostForm("username"))
	realname := strings.TrimSpace(c.PostForm("realname"))
	if realname == "" {
		realname = username
	}
	var exist string
	if username == "" {
		message = "Username is required."
	} else if err := db.QueryRow("SELECT id FROM user WHERE username = ? AND id != ?", username, id).Scan(&exist); err == nil {
		message = fmt.Sprintf("Username %s is already existed.", username)
	} else if deptID == "" {
		message = "Department is required."
	} else {
		if _, err = db.Exec("UPDATE user SET username = ?, realname = ?, dept_id = ?, type = ?, permission = ? WHERE id = ?",
			username, realname, deptID, id); err != nil {
			log.Printf("Failed to edit user: %v", err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		return
	}
	c.JSON(200, gin.H{"status": 0, "message": message})
}

func doDeleteEmpl(c *gin.Context) {
	id := c.Param("id")
	empl, err := getEmpl(id)
	if err != nil {
		log.Printf("Failed to get empl: %v", err)
		c.String(400, "")
		return
	}
	if !checkPermission(strconv.Itoa(empl.DeptID), c) {
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
