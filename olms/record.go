package olms

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type record struct {
	ID       int
	DeptID   int
	DeptName string
	UserID   int
	Name     string
	Date     time.Time
	Duration int
	Type     bool
	Status   int
	Describe string
	Created  time.Time
}

func getRecords(id interface{}, deptIDs []string, year, month, Type, status string, page interface{}) (records []record, total int, err error) {
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	stmt := "SELECT %s FROM record JOIN employee ON user_id = employee.id WHERE "

	var args []interface{}
	if year != "" {
		if month == "" {
			stmt += "strftime('%Y', date) = ? AND "
			args = append(args, year)
		}
	} else {
		stmt += "strftime('%Y%m', date) = ? AND "
		args = append(args, year+month)
	}
	if Type != "" {
		stmt += "type = ? AND "
		args = append(args, Type)
	}
	if status != "" {
		stmt += "status = ? AND "
		args = append(args, Type)
	}

	if id != nil {
		stmt += " user_id = ?"
		args = append(args, id)
	} else {
		marks := make([]string, len(deptIDs))
		for i := range marks {
			marks[i] = "?"
		}
		stmt += " record.dept_id IN (" + strings.Join(marks, ", ") + ")"
		for _, i := range deptIDs {
			args = append(args, i)
		}
	}

	var limit string
	bc := make(chan bool, 1)
	if p, ok := page.(int); ok {
		limit = fmt.Sprintf(" LIMIT ?, ?")
		args = append(args, (p-1)*perPage, perPage)
		go func() {
			if err := db.QueryRow(fmt.Sprintf(stmt, "count(*)")).Scan(&total); err != nil {
				log.Printf("Failed to get total records: %v", err)
				bc <- false
			}
			bc <- true
		}()
	} else {
		bc <- true
	}
	rows, err := db.Query(fmt.Sprintf(stmt+" ORDER BY created DESC"+limit,
		"record.id, employee.dept_id, dept_name, employee.user_id, realname, date, ABS(duration), type, status, describe, created"), args...)
	if err != nil {
		log.Printf("Failed to get records: %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r record
		if err = rows.Scan(&r.ID, &r.DeptID, &r.DeptName, &r.UserID, &r.Name, &r.Date, &r.Duration, &r.Type, &r.Status, &r.Describe, &r.Created); err != nil {
			log.Printf("Failed to scan record: %v", err)
			return
		}
		records = append(records, r)
	}
	if v := <-bc; !v {
		err = fmt.Errorf("Failed to get total records")
	}
	return
}

func checkRecord(c *gin.Context, record record, super bool) bool {
	session := sessions.Default(c)
	users, _, err := getEmpls(session.Get("userID"), nil, nil, nil)
	if err != nil {
		return false
	}
	if super {
		if record.UserID == users[0].ID || users[0].ID == 0 {
			return true
		}
		return false
	}
	for _, i := range strings.Split(users[0].Permission, ",") {
		if strconv.Itoa(record.DeptID) == i {
			return true
		}
	}
	return false
}

func showRecords(c *gin.Context) {
	c.HTML(200, "showRecords.html", nil)
}

func addRecord(c *gin.Context) {
	c.HTML(200, "addRecord.html", nil)
}

func doAddRecord(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	date := c.PostForm("date")
	Type, err := strconv.Atoi(c.PostForm("type"))
	if err != nil {
		log.Printf("Failed to get type: %v", err)
		c.String(400, "")
		return
	}
	duration, err := strconv.Atoi(c.PostForm("duration"))
	if err != nil {
		log.Printf("Failed to get duration: %v", err)
		c.String(400, "")
		return
	}
	switch Type {
	case 1:
		if duration < 1 {
			c.JSON(200, gin.H{"status": 0, "message": "Bad duration value"})
			return
		}
	case 0:
		duration = -duration
		if duration > -1 {
			c.JSON(200, gin.H{"status": 0, "message": "Bad duration value"})
			return
		}
	default:
		log.Println("Unknown type value")
		c.String(400, "")
		return
	}
	describe := c.PostForm("describe")

	session := sessions.Default(c)
	users, _, err := getEmpls(session.Get("userID"), nil, nil, nil)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		c.String(500, "")
		return
	}
	userID := c.PostForm("empl")
	ip := c.ClientIP()
	if userID == "" {
		if _, err := db.Exec("INSERT INTO record (date, type, duration, describe, dept_id, user_id, createdby) VALUES (?, ?, ?, ?, ?, ?, ?)",
			date, Type, duration, describe, users[0].DeptID, users[0].ID, fmt.Sprintf("%d-%s", users[0].ID, ip)); err != nil {
			log.Printf("Failed to add record: %v", err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		return
	}
	deptID := c.PostForm("dept")
	if deptID != "" && !checkPermission(c, deptID, userID) {
		c.String(403, "")
		return
	}
	if _, err := db.Exec("INSERT INTO record (dept_id, user_id, date, type, duration, describe, status, createdby, verifiedby) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		deptID, userID, date, Type, duration, describe, 1, fmt.Sprintf("%d-%s", users[0].ID, ip), fmt.Sprintf("%d-%s", users[0].ID, ip)); err != nil {
		log.Printf("Failed to add record: %v", err)
		c.String(500, "")
		return
	}
	c.JSON(200, gin.H{"status": 1})
}

func editRecord(c *gin.Context) {
	id := c.Param("id")
	records, _, err := getRecords(id, nil, "", "", "", "", nil)
	if err != nil {
		log.Printf("Failed to get record: %v", err)
		c.String(400, "")
		return
	}
	if !checkRecord(c, records[0], true) {
		c.String(403, "")
		return
	}
	c.HTML(200, "editRecord.html", gin.H{"record": records[0]})
}

func doEditRecord(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	date := c.PostForm("date")
	Type, err := strconv.Atoi(c.PostForm("type"))
	if err != nil {
		log.Printf("Failed to get type: %v", err)
		c.String(400, "")
		return
	}
	duration, err := strconv.Atoi(c.PostForm("duration"))
	if err != nil {
		log.Printf("Failed to get duration: %v", err)
		c.String(400, "")
		return
	}
	switch Type {
	case 1:
		if duration < 1 {
			c.JSON(200, gin.H{"status": 0, "message": "Bad duration value"})
			return
		}
	case 0:
		duration = -duration
		if duration > -1 {
			c.JSON(200, gin.H{"status": 0, "message": "Bad duration value"})
			return
		}
	default:
		log.Println("Unknown type value")
		c.String(400, "")
		return
	}
	describe := c.PostForm("describe")

	id := c.Param("id")
	records, _, err := getRecords(id, nil, "", "", "", "", nil)
	if err != nil {
		log.Printf("Failed to get record: %v", err)
		c.String(400, "")
		return
	}
	if !checkRecord(c, records[0], true) {
		c.String(403, "")
		return
	}
	userID := c.PostForm("empl")
	if userID == "" {
		if records[0].Status != 0 {
			c.JSON(200, gin.H{"status": 0, "message": "You can only update record which is not verified."})
			return
		}
		if _, err := db.Exec("UPDATE record SET date = ?, type = ?, duration = ?, describe = ? WHERE id = ?",
			date, Type, duration, describe, id); err != nil {
			log.Printf("Failed to update record: %v", err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		return
	}
	deptID := c.PostForm("dept")
	if deptID != "" && !checkPermission(c, deptID, userID) {
		c.String(403, "")
		return
	}
	status := c.PostForm("status")
	if _, err := db.Exec("UPDATE record SET empl_id = ?, dept_id = ?, date = ?, type = ?, duration = ?, status = ?, describe = ? WHERE id = ?",
		userID, deptID, date, Type, duration, status, describe, id); err != nil {
		log.Printf("Failed to update record: %v", err)
		c.String(500, "")
		return
	}
	c.JSON(200, gin.H{"status": 1})
}

func verifyRecord(c *gin.Context) {
	id := c.Param("id")
	records, _, err := getRecords(id, nil, "", "", "", "", nil)
	if err != nil {
		log.Printf("Failed to get record: %v", err)
		c.String(400, "")
		return
	}
	if !checkRecord(c, records[0], false) {
		c.String(403, "")
		return
	}
	c.HTML(200, "verifyRecord.html", gin.H{"record": records[0]})
}

func doVerifyRecord(c *gin.Context) {
	id := c.Param("id")
	records, _, err := getRecords(id, nil, "", "", "", "", nil)
	if err != nil {
		log.Printf("Failed to get record: %v", err)
		c.String(400, "")
		return
	}
	if !checkRecord(c, records[0], false) {
		c.String(403, "")
		return
	}
	if records[0].Status != 0 {
		log.Println("The record is already verified.")
		c.String(400, "")
		return
	}
	status, err := strconv.Atoi(c.PostForm("status"))
	if err != nil {
		log.Printf("Failed to get status: %v", err)
		c.String(400, "")
		return
	}
	if status != 1 && status != 2 {
		log.Printf("Unknow status.")
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

	session := sessions.Default(c)
	users, _, err := getEmpls(session.Get("userID"), nil, nil, nil)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		c.String(500, "")
		return
	}
	ip := c.ClientIP()
	if _, err := db.Exec("UPDATE record SET status = ?, verifiedby = ? WHERE id = ?", status, fmt.Sprintf("%d-%s", users[0].ID, ip), id); err != nil {
		log.Printf("Failed to verify record: %v", err)
		c.String(500, "")
		return
	}
	c.JSON(200, gin.H{"status": 1})
}

func doDeleteRecord(c *gin.Context) {
	id := c.Param("id")
	records, _, err := getRecords(id, nil, "", "", "", "", nil)
	if err != nil {
		log.Printf("Failed to get record: %v", err)
		c.String(400, "")
		return
	}
	session := sessions.Default(c)
	users, _, err := getEmpls(session.Get("userID"), nil, nil, nil)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		c.String(500, "")
		return
	}
	if users[0].ID != 0 && records[0].Status != 0 {
		c.JSON(200, gin.H{"status": 0, "message": "You can only delete record which is not verified."})
		return
	}

	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	if _, err := db.Exec("DELETE FROM record WHERE id = ?", id); err != nil {
		log.Printf("Failed to delete record: %v", err)
		c.String(500, "")
		return
	}
	c.JSON(200, gin.H{"status": 1})
}
