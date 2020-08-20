package olms

import (
	"fmt"
	"log"
	"net"
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

func getRecords(id *idOptions, options *searchOptions) (records []record, total int, err error) {
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	stmt := "SELECT %s FROM record JOIN employee ON user_id = employee.id WHERE"

	var args []interface{}
	var orderBy, limit string
	bc := make(chan bool, 1)
	if id.RecordID != nil {
		stmt += " record.id = ?"
		args = append(args, id.RecordID)
		bc <- true
	} else {
		if id.UserID != nil {
			stmt += " user_id = ?"
			args = append(args, id.UserID)
		} else {
			marks := make([]string, len(id.DeptIDs))
			for i := range marks {
				marks[i] = "?"
			}
			stmt += " record.dept_id IN (" + strings.Join(marks, ", ") + ")"
			for _, i := range id.DeptIDs {
				args = append(args, i)
			}
		}

		if options != nil {
			if options.Year != nil {
				if options.Month == nil {
					stmt += " AND strftime('%%Y', date) = ?"
					args = append(args, options.Year)
				} else {
					stmt += " AND strftime('%%Y%%m', date) = ?"
					args = append(args, fmt.Sprintf("%v%v", options.Year, options.Month))
				}
			}
			if options.Type != nil {
				stmt += " AND type = ?"
				args = append(args, options.Type)
			}
			if options.Status != nil {
				stmt += " AND status = ?"
				args = append(args, options.Status)
			}
			if options.Describe != nil {
				stmt += " AND describe LIKE ?"
				args = append(args, fmt.Sprintf("%%%v%%", options.Describe))
			}
			if p, ok := options.Page.(float64); ok {
				limit = fmt.Sprintf(" LIMIT ?, ?")
				args = append(args, int(p-1)*perPage, perPage)
			} else {
				bc <- true
			}
			if options.Sort != nil {
				orderBy = fmt.Sprintf(" ORDER BY %v %v", options.Sort, options.Order)
			} else {
				orderBy = " ORDER BY created DESC"
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
	rows, err := db.Query(fmt.Sprintf(stmt+orderBy+limit,
		"record.id, employee.dept_id, dept_name, employee.id, realname, date, ABS(duration), type, status, describe, created"), args...)
	if err != nil {
		log.Printf("Failed to get records: %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r record
		if err = rows.Scan(
			&r.ID, &r.DeptID, &r.DeptName, &r.UserID, &r.Name, &r.Date, &r.Duration, &r.Type, &r.Status, &r.Describe, &r.Created); err != nil {
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

func getYears(id *idOptions) (years []string, err error) {
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	stmt := "SELECT DISTINCT strftime('%Y', date) year FROM record WHERE"

	var args []interface{}
	if id.UserID != nil {
		stmt += " user_id = ?"
		args = append(args, id.UserID)
	} else {
		marks := make([]string, len(id.DeptIDs))
		for i := range marks {
			marks[i] = "?"
		}
		stmt += " dept_id IN (" + strings.Join(marks, ", ") + ")"
		for _, i := range id.DeptIDs {
			args = append(args, i)
		}
	}
	rows, err := db.Query(stmt+" ORDER BY year DESC", args...)
	if err != nil {
		log.Printf("Failed to get years: %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var y string
		if err = rows.Scan(&y); err != nil {
			log.Printf("Failed to scan year: %v", err)
			return
		}
		years = append(years, y)
	}
	return
}

func checkRecord(c *gin.Context, record record, super bool) bool {
	userID := sessions.Default(c).Get("userID")
	if userID == "0" {
		return true
	}
	users, _, err := getEmpls(&idOptions{UserID: userID}, nil)
	if err != nil {
		return false
	}
	if super {
		if record.UserID == users[0].ID {
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

func doAddRecord(c *gin.Context) {
	if !verifyResponse("record", c.ClientIP(), c.PostForm("g-recaptcha-response")) {
		c.String(403, "reCAPTCHA challenge failed")
		return
	}
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

	var user empl
	switch userID := sessions.Default(c).Get("userID"); userID {
	case "0":
		user = empl{ID: 0}
	default:
		users, _, err := getEmpls(&idOptions{UserID: userID}, nil)
		if err != nil {
			log.Printf("Failed to get user: %v", err)
			c.String(500, "")
			return
		}
		user = users[0]
	}

	userID := c.PostForm("empl")
	ip, _, _ := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr))
	ip = ip + "-" + c.ClientIP()
	if userID == "" {
		if _, err := db.Exec("INSERT INTO record (date, type, duration, describe, dept_id, user_id, createdby) VALUES (?, ?, ?, ?, ?, ?, ?)",
			date, Type, duration, describe, user.DeptID, user.ID, fmt.Sprintf("%d-%s", user.ID, ip)); err != nil {
			log.Printf("Failed to add record: %v", err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		return
	}
	deptID := c.PostForm("dept")
	if user.ID != 0 {
		if deptID != "" && !checkPermission(c, deptID, userID) {
			c.String(403, "")
			return
		}
		if _, err := db.Exec("INSERT INTO record (dept_id, user_id, date, type, duration, describe, status, createdby, verifiedby) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
			deptID, userID, date, Type, duration, describe, 1, fmt.Sprintf("%d-%s", user.ID, ip), fmt.Sprintf("%d-%s", user.ID, ip)); err != nil {
			log.Printf("Failed to add record: %v", err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		return
	}
	status := c.PostForm("status")
	if _, err := db.Exec("INSERT INTO record (dept_id, user_id, date, type, duration, describe, status, createdby, verifiedby) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		deptID, userID, date, Type, duration, describe, status, fmt.Sprintf("%d-%s", 0, ip), fmt.Sprintf("%d-%s", 0, ip)); err != nil {
		log.Printf("Failed to add record: %v", err)
		c.String(500, "")
		return
	}
	c.JSON(200, gin.H{"status": 1})
}

func doEditRecord(c *gin.Context) {
	if !verifyResponse("record", c.ClientIP(), c.PostForm("g-recaptcha-response")) {
		c.String(403, "reCAPTCHA challenge failed")
		return
	}
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
	records, _, err := getRecords(&idOptions{RecordID: id}, nil)
	if err != nil {
		log.Printf("Failed to get record: %v", err)
		c.String(400, "")
		return
	}
	if !checkRecord(c, records[0], true) {
		c.String(403, "")
		return
	}
	user := sessions.Default(c).Get("userID")
	if user != "0" {
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
	userID := c.PostForm("empl")
	deptID := c.PostForm("dept")
	if userID == "" || deptID == "" {
		log.Print("Missing param.")
		c.String(400, "")
		return
	}
	users, _, err := getEmpls(&idOptions{UserID: userID}, nil)
	if err != nil {
		log.Printf("Failed to get users: %v", err)
		c.String(400, "")
		return
	}
	if deptID != strconv.Itoa(users[0].DeptID) {
		c.JSON(200, gin.H{"status": 0, "message": "Employee does not belong this department."})
		return
	}
	status := c.PostForm("status")
	if _, err := db.Exec("UPDATE record SET user_id = ?, dept_id = ?, date = ?, type = ?, duration = ?, status = ?, describe = ? WHERE id = ?",
		userID, deptID, date, Type, duration, status, describe, id); err != nil {
		log.Printf("Failed to update record: %v", err)
		c.String(500, "")
		return
	}
	c.JSON(200, gin.H{"status": 1})
}

func verifyRecord(c *gin.Context) {
	id := c.Param("id")
	records, _, err := getRecords(&idOptions{RecordID: id}, nil)
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
	if !verifyResponse("verify", c.ClientIP(), c.PostForm("g-recaptcha-response")) {
		c.String(403, "reCAPTCHA challenge failed")
		return
	}
	id := c.Param("id")
	records, _, err := getRecords(&idOptions{RecordID: id}, nil)
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

	var user empl
	switch userID := sessions.Default(c).Get("userID"); userID {
	case "0":
		user = empl{ID: 0}
	default:
		users, _, err := getEmpls(&idOptions{UserID: userID}, nil)
		if err != nil {
			log.Printf("Failed to get user: %v", err)
			c.String(500, "")
			return
		}
		user = users[0]
	}
	ip, _, _ := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr))
	if _, err := db.Exec("UPDATE record SET status = ?, verifiedby = ? WHERE id = ?",
		status, fmt.Sprintf("%d-%s-%s", user.ID, ip, c.ClientIP()), id); err != nil {
		log.Printf("Failed to verify record: %v", err)
		c.String(500, "")
		return
	}
	c.JSON(200, gin.H{"status": 1})
}

func doDeleteRecord(c *gin.Context) {
	id := c.Param("id")
	records, _, err := getRecords(&idOptions{RecordID: id}, nil)
	if err != nil {
		log.Printf("Failed to get record: %v", err)
		c.String(400, "")
		return
	}
	var user empl
	switch userID := sessions.Default(c).Get("userID"); userID {
	case "0":
		user = empl{ID: 0}
	default:
		users, _, err := getEmpls(&idOptions{UserID: userID}, nil)
		if err != nil {
			log.Printf("Failed to get user: %v", err)
			c.String(500, "")
			return
		}
		user = users[0]
	}
	if user.ID != 0 && records[0].Status != 0 {
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
