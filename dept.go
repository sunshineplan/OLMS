package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type dept struct {
	ID   int
	Name string
}

func getDept(id interface{}) (dept, error) {
	var dept dept
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return dept, err
	}
	defer db.Close()
	if err := db.QueryRow("SELECT * FROM department WHERE id = ?", id).Scan(&dept.ID, &dept.Name); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Printf("No results: %v", err)
		} else {
			log.Printf("Failed to query dept: %v", err)
		}
		return dept, err
	}
	return dept, nil
}

func getDepts(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.String(503, "")
		return
	}
	defer db.Close()
	session := sessions.Default(c)
	user, err := getEmpl(session.Get("userID"))
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		c.String(503, "")
		return
	}

	var depts []dept
	rows, err := db.Query("SELECT * FROM department WHERE id IN (?) ORDER BY dept_name", user.Permission)
	if err != nil {
		log.Printf("Failed to get departments: %v", err)
		c.String(500, "")
		return
	}
	defer rows.Close()
	for rows.Next() {
		var dept dept
		if err := rows.Scan(&dept.ID, &dept.Name); err != nil {
			log.Printf("Failed to scan department: %v", err)
			c.String(500, "")
			return
		}
		depts = append(depts, dept)
	}
	c.JSON(200, depts)
}

func showDept(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func addDept(c *gin.Context) {
	c.HTML(200, "addDept.html", nil)
}

func doAddDept(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.String(503, "")
		return
	}
	defer db.Close()
	dept := strings.TrimSpace(c.PostForm("dept"))
	var exist, message string
	if dept == "" {
		message = "Department name is required."
	} else if err := db.QueryRow("SELECT id FROM department WHERE dept_name = ?", dept).Scan(&exist); err == nil {
		message = fmt.Sprintf("Department %s is already existed.", dept)
	} else {
		if _, err := db.Exec("INSERT INTO department (dept_name) VALUES (?)", dept); err != nil {
			log.Printf("Failed to add department: %v", err)
			c.String(500, "")
			return
		}
		if _, err := db.Exec("UPDATE employee SET permission = (SELECT group_concat(id) FROM department) WHERE id = 0"); err != nil {
			log.Printf("Failed to add admin permission: %v", err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		return
	}
	c.JSON(200, gin.H{"status": 0, "message": message})
}

func editDept(c *gin.Context) {
	dept, err := getDept(c.Param("id"))
	if err != nil {
		log.Printf("Failed to get dept id: %v", err)
		c.String(400, "")
		return
	}
	c.HTML(200, "addDept.html", gin.H{"dept": dept})
}

func doEditDept(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.String(503, "")
		return
	}
	defer db.Close()
	dept := strings.TrimSpace(c.PostForm("dept"))
	id := c.Param("id")
	var exist, message string
	if dept == "" {
		message = "Department name is required."
	} else if err := db.QueryRow("SELECT id FROM department WHERE dept_name = ? AND id != ?", id, dept).Scan(&exist); err == nil {
		message = fmt.Sprintf("Department %s is already existed.", dept)
	} else {
		if _, err := db.Exec("UPDATE department SET dept_name = ? WHERE id = ?", dept, id); err != nil {
			log.Printf("Failed to edit department: %v", err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		return
	}
	c.JSON(200, gin.H{"status": 0, "message": message})
}

func doDeleteDept(c *gin.Context) {
	id := c.Param("id")
	if _, err := getDept(id); err != nil {
		log.Printf("Failed to get dept id: %v", err)
		c.String(400, "")
		return
	}
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.String(503, "")
		return
	}
	defer db.Close()
	if _, err := db.Exec("DELETE FROM department WHERE id = ?", id); err != nil {
		log.Printf("Failed to delete department: %v", err)
		c.String(500, "")
		return
	}
	c.JSON(200, gin.H{"status": 1})
}
