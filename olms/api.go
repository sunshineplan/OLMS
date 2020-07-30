package olms

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

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

func getEmpls(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	stmt := "SELECT %s FROM employee WHERE "

	deptID := c.Query("dept")
	if deptID != "" && !checkPermission(deptID, c) {
		c.String(403, "")
		return
	}
	if deptID != "" {
		stmt += "dept_id = " + deptID
	} else {
		session := sessions.Default(c)
		user, err := getEmpl(session.Get("userID"))
		if err != nil {
			log.Printf("Failed to connect to database: %v", err)
			c.String(500, "")
			return
		}
		stmt += fmt.Sprintf("dept_id IN (%s)", user.Permission)
	}
	if Type := c.Query("type"); Type != "" {
		stmt += " AND type = " + Type
	}

	bc := make(chan bool, 1)
	var records int
	go func() {
		if err := db.QueryRow(fmt.Sprintf(stmt, "count(*)")).Scan(&records); err != nil {
			log.Printf("Failed to get records: %v", err)
			bc <- false
		}
		bc <- true
	}()

	stmt += " ORDER BY dept_name, realname"
	if page, err := strconv.Atoi(c.Query("page")); err != nil {
		stmt += fmt.Sprintf(" LIMIT %d, %d", (page-1)*perPage, perPage)
	}
	var empls []empl
	rows, err := db.Query(fmt.Sprintf(stmt, "id, realname, dept_name"))
	if err != nil {
		log.Printf("Failed to get employees: %v", err)
		c.String(500, "")
		return
	}
	defer rows.Close()
	for rows.Next() {
		var empl empl
		if err := rows.Scan(&empl.ID, &empl.Realname, &empl.DeptName); err != nil {
			log.Printf("Failed to scan employee: %v", err)
			c.String(500, "")
			return
		}
		empls = append(empls, empl)
	}
	if v := <-bc; !v {
		c.String(500, "")
		return
	}
	c.JSON(200, gin.H{"empls": empls, "records": records})
}
