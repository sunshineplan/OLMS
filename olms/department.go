package olms

import (
	"database/sql"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

type department struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func getDepartments(db *sql.DB, ids []string, super bool) (departments []department, err error) {
	var rows *sql.Rows
	if super {
		rows, err = db.Query("SELECT * FROM department")
		if err != nil {
			log.Println("Failed to get departments:", err)
			return
		}
	} else {
		rows, err = db.Query("SELECT * FROM department WHERE id IN (" + strings.Join(ids, ", ") + ")")
		if err != nil {
			log.Println("Failed to get departments:", err)
			return
		}
	}
	defer rows.Close()

	for rows.Next() {
		var department department
		if err = rows.Scan(&department.ID, &department.Name); err != nil {
			log.Println("Failed to scan department:", err)
			return
		}
		departments = append(departments, department)
	}
	return
}

func addDepartment(c *gin.Context) {
	var department department
	if err := c.BindJSON(&department); err != nil {
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
	var exist, message string
	if department.Name == "" {
		message = "DepartmentRequired"
	} else if err := db.QueryRow("SELECT id FROM department WHERE dept_name = ?", department.Name).Scan(&exist); err == nil {
		message = "DepartmentExist"
	} else {
		if _, err := db.Exec("INSERT INTO department (dept_name) VALUES (?)", department.Name); err != nil {
			log.Println("Failed to add department:", err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		return
	}
	c.JSON(200, gin.H{"status": 0, "message": message})
}

func editDepartment(c *gin.Context) {
	var department department
	if err := c.BindJSON(&department); err != nil {
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
	var old, exist, message string
	if department.Name == "" {
		message = "DepartmentRequired"
	} else if db.QueryRow("SELECT dept_name FROM department WHERE id = ?", department.ID).Scan(&old); old == department.Name {
		message = "SameDepartment"
	} else if err := db.QueryRow("SELECT id FROM department WHERE dept_name = ? AND id != ?",
		department.ID, department.Name).Scan(&exist); err == nil {
		message = "DepartmentExist"
	} else {
		if _, err := db.Exec("UPDATE department SET dept_name = ? WHERE id = ?",
			department.Name, department.ID); err != nil {
			log.Println("Failed to edit department:", err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		return
	}
	c.JSON(200, gin.H{"status": 0, "message": message})
}

func deleteDepartment(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	id := c.Param("id")
	if _, err := db.Exec("DELETE FROM department WHERE id = ?", id); err != nil {
		log.Println("Failed to delete department:", err)
		c.String(500, "")
		return
	}
	c.JSON(200, gin.H{"status": 1})
}
