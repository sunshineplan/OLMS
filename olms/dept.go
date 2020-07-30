package olms

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

type dept struct {
	ID   int
	Name string
}

func showDept(c *gin.Context) {
	c.HTML(200, "showDept.html", nil)
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
