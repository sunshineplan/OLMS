package olms

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

type department struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func getDepartments(ids []string) ([]department, error) {
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, err
	}
	defer db.Close()

	marks := make([]string, len(ids))
	for i := range marks {
		marks[i] = "?"
	}
	var departments []department
	var args []interface{}
	for _, i := range ids {
		args = append(args, i)
	}
	rows, err := db.Query("SELECT * FROM department WHERE id IN ("+strings.Join(marks, ", ")+")", args...)
	if err != nil {
		log.Printf("Failed to get departments: %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var department department
		if err := rows.Scan(&department.ID, &department.Name); err != nil {
			log.Printf("Failed to scan department: %v", err)
			return nil, err
		}
		departments = append(departments, department)
	}
	return departments, nil
}

func addDepartment(c *gin.Context) {
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
		message = fmt.Sprintf(localize(c)["DepartmentExist"], dept)
	} else {
		if _, err := db.Exec("INSERT INTO department (dept_name) VALUES (?)", dept); err != nil {
			log.Printf("Failed to add department: %v", err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		return
	}
	c.JSON(200, gin.H{"status": 0, "message": message})
}

func editDepartment(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.String(503, "")
		return
	}
	defer db.Close()
	dept := strings.TrimSpace(c.PostForm("dept"))
	id := c.Param("id")
	localize := localize(c)
	var old, exist, message string
	if dept == "" {
		message = "Department name is required."
	} else if db.QueryRow("SELECT dept_name FROM department WHERE id = ?", id).Scan(&old); old == dept {
		message = localize["SameDepartment"]
	} else if err := db.QueryRow("SELECT id FROM department WHERE dept_name = ? AND id != ?", id, dept).Scan(&exist); err == nil {
		message = fmt.Sprintf(localize["DepartmentExist"], dept)
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

func deleteDepartment(c *gin.Context) {
	id := c.Param("id")
	if _, err := getDepartments([]string{id}); err != nil {
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
