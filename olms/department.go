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

func getDepartments(db *sql.DB, ids []string, super bool) ([]department, error) {
	var rows *sql.Rows
	var err error
	if super {
		rows, err = db.Query("SELECT * FROM department")
		if err != nil {
			log.Println("Failed to get departments:", err)
			return nil, err
		}
	} else {
		rows, err = db.Query("SELECT * FROM department WHERE id IN (" + strings.Join(ids, ", ") + ")")
		if err != nil {
			log.Println("Failed to get departments:", err)
			return nil, err
		}
	}
	defer rows.Close()

	departments := []department{}
	for rows.Next() {
		var department department
		if err := rows.Scan(&department.ID, &department.Name); err != nil {
			log.Println("Failed to scan department:", err)
			return nil, err
		}
		departments = append(departments, department)
	}
	return departments, nil
}

func addDepartment(c *gin.Context) {
	var department department
	if err := c.BindJSON(&department); err != nil {
		log.Println("Failed to get option:", err)
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
	} else if err := db.QueryRow("SELECT id FROM department WHERE deptname = ?", department.Name).Scan(&exist); err == nil {
		message = "DepartmentExist"
	} else {
		if _, err := db.Exec("INSERT INTO department (deptname) VALUES (?)", department.Name); err != nil {
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
		log.Println("Failed to get option:", err)
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
	} else if db.QueryRow("SELECT deptname FROM department WHERE id = ?", department.ID).Scan(&old); old == department.Name {
		message = "SameDepartment"
	} else if err := db.QueryRow("SELECT id FROM department WHERE deptname = ? AND id != ?",
		department.ID, department.Name).Scan(&exist); err == nil {
		message = "DepartmentExist"
	} else {
		if _, err := db.Exec("UPDATE department SET deptname = ? WHERE id = ?",
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
	var exist string
	if err := db.QueryRow("SELECT realname FROM user WHERE dept_id = ?", id).Scan(&exist); err == nil {
		c.String(500, "DeleteNonEmptyDepartment")
		return
	}
	if _, err := db.Exec("DELETE FROM department WHERE id = ?", id); err != nil {
		log.Println("Failed to delete department:", err)
		c.String(500, "")
		return
	}
	c.JSON(200, gin.H{"status": 1})
}
