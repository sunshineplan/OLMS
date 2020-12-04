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
	ID       int       `json:"id"`
	DeptID   int       `json:"deptid"`
	DeptName string    `json:"deptname"`
	UserID   int       `json:"userid"`
	Realname string    `json:"realname"`
	Date     time.Time `json:"date"`
	Duration int       `json:"duration"`
	Type     bool      `json:"type"`
	Status   int       `json:"status"`
	Describe string    `json:"describe"`
	Created  time.Time `json:"created"`
}

func getRecords(id *idOptions, options *searchOptions) (records []record, total int, err error) {
	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		return
	}
	defer db.Close()

	stmt := "SELECT %s FROM record JOIN employee ON user_id = employee.id WHERE"

	var args []interface{}
	var orderBy, limit string
	bc := make(chan bool, 1)
	if id.Record != nil {
		stmt += " record.id = ?"
		args = append(args, id.Record)
		bc <- true
	} else {
		if id.User != nil {
			stmt += " user_id = ?"
			args = append(args, id.User)
		} else {
			marks := make([]string, len(id.Departments))
			for i := range marks {
				marks[i] = "?"
			}
			stmt += " record.dept_id IN (" + strings.Join(marks, ", ") + ")"
			for _, i := range id.Departments {
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
				log.Println("Failed to get total records:", err)
				bc <- false
			}
			bc <- true
		}()
	}
	rows, err := db.Query(fmt.Sprintf(stmt+orderBy+limit,
		"record.id, employee.dept_id, dept_name, employee.id, realname, date, ABS(duration), type, status, describe, created"), args...)
	if err != nil {
		log.Println("Failed to get records:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r record
		if err = rows.Scan(
			&r.ID, &r.DeptID, &r.DeptName, &r.UserID, &r.Realname, &r.Date, &r.Duration, &r.Type, &r.Status, &r.Describe, &r.Created,
		); err != nil {
			log.Println("Failed to scan record:", err)
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
		log.Println("Failed to connect to database:", err)
		return
	}
	defer db.Close()

	stmt := "SELECT DISTINCT strftime('%Y', date) year FROM record WHERE"

	var args []interface{}
	if id.User != nil {
		stmt += " user_id = ?"
		args = append(args, id.User)
	} else {
		marks := make([]string, len(id.Departments))
		for i := range marks {
			marks[i] = "?"
		}
		stmt += " dept_id IN (" + strings.Join(marks, ", ") + ")"
		for _, i := range id.Departments {
			args = append(args, i)
		}
	}
	rows, err := db.Query(stmt+" ORDER BY year DESC", args...)
	if err != nil {
		log.Println("Failed to get years:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var y string
		if err = rows.Scan(&y); err != nil {
			log.Println("Failed to scan year:", err)
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
	user, err := getUser(userID)
	if err != nil {
		return false
	}
	if super {
		if record.UserID == user.ID {
			return true
		}
		return false
	}
	for _, i := range strings.Split(user.Permission, ",") {
		if strconv.Itoa(record.DeptID) == i {
			return true
		}
	}
	return false
}

func addRecord(c *gin.Context) {
	if !verifyResponse("record", c.ClientIP(), c.PostForm("g-recaptcha-response")) {
		c.String(403, "reCAPTCHA challenge failed")
		return
	}
	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	date := c.PostForm("date")
	Type, err := strconv.Atoi(c.PostForm("type"))
	if err != nil {
		log.Println("Failed to get type:", err)
		c.String(400, "")
		return
	}
	duration, err := strconv.Atoi(c.PostForm("duration"))
	if err != nil {
		log.Println("Failed to get duration:", err)
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

	var user employee
	switch userID := sessions.Default(c).Get("userID"); userID {
	case "0":
		user = employee{ID: 0}
	default:
		user, err = getUser(userID)
		if err != nil {
			log.Println("Failed to get user:", err)
			c.String(500, "")
			return
		}
	}

	localize := localize(c)
	userID := c.PostForm("empl")
	ip, _, _ := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr))
	ip = ip + "-" + c.ClientIP()
	if userID == "" {
		if _, err := db.Exec(
			"INSERT INTO record (date, type, duration, describe, dept_id, user_id, createdby) VALUES (?, ?, ?, ?, ?, ?, ?)",
			date, Type, duration, describe, user.DeptID, user.ID, fmt.Sprintf("%d-%s", user.ID, ip)); err != nil {
			log.Println("Failed to add record:", err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		notify(&idOptions{Departments: []string{strconv.Itoa(user.DeptID)}},
			fmt.Sprintf(localize["AddRecordSubscribe"], user.Realname), localize)
		return
	}
	deptID := c.PostForm("dept")
	if user.ID != 0 {
		if deptID != "" && !checkPermission(c, &idOptions{User: userID, Departments: []string{fmt.Sprintf("%v", deptID)}}) {
			c.String(403, "")
			return
		}
		if _, err := db.Exec(
			"INSERT INTO record (dept_id, user_id, date, type, duration, describe, status, createdby, verifiedby) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
			deptID, userID, date, Type, duration, describe, 1, fmt.Sprintf("%d-%s", user.ID, ip), fmt.Sprintf("%d-%s", user.ID, ip),
		); err != nil {
			log.Println("Failed to add record:", err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		notify(&idOptions{User: userID},
			fmt.Sprintf(localize["AdminAddRecordSubscribe"], user.Realname), localize)
		notify(&idOptions{Departments: []string{deptID}},
			fmt.Sprintf(localize["AdminAddRecordAdminSubscribe"], user.Realname), localize)
		return
	}
	status := c.PostForm("status")
	if _, err := db.Exec(
		"INSERT INTO record (dept_id, user_id, date, type, duration, describe, status, createdby, verifiedby) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		deptID, userID, date, Type, duration, describe, status, fmt.Sprintf("%d-%s", 0, ip), fmt.Sprintf("%d-%s", 0, ip)); err != nil {
		log.Println("Failed to add record:", err)
		c.String(500, "")
		return
	}
	c.JSON(200, gin.H{"status": 1})
}

func editRecord(c *gin.Context) {
	if !verifyResponse("record", c.ClientIP(), c.PostForm("g-recaptcha-response")) {
		c.String(403, "reCAPTCHA challenge failed")
		return
	}
	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	date := c.PostForm("date")
	Type, err := strconv.Atoi(c.PostForm("type"))
	if err != nil {
		log.Println("Failed to get type:", err)
		c.String(400, "")
		return
	}
	duration, err := strconv.Atoi(c.PostForm("duration"))
	if err != nil {
		log.Println("Failed to get duration:", err)
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
	records, _, err := getRecords(&idOptions{Record: id}, nil)
	if err != nil {
		log.Println("Failed to get record:", err)
		c.String(400, "")
		return
	}
	if !checkRecord(c, records[0], true) {
		c.String(403, "")
		return
	}
	userID := sessions.Default(c).Get("userID")
	if userID != "0" {
		if records[0].Status != 0 {
			c.JSON(200, gin.H{"status": 0, "message": "You can only update record which is not verified."})
			return
		}
		if _, err := db.Exec("UPDATE record SET date = ?, type = ?, duration = ?, describe = ? WHERE id = ?",
			date, Type, duration, describe, id); err != nil {
			log.Println("Failed to update record:", err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		localize := localize(c)
		notify(&idOptions{Departments: []string{strconv.Itoa(records[0].DeptID)}},
			fmt.Sprintf(localize["EditRecordSubscribe"], records[0].Realname), localize)
		return
	}
	emplID := c.PostForm("empl")
	deptID := c.PostForm("dept")
	if emplID == "" || deptID == "" {
		log.Print("Missing param.")
		c.String(400, "")
		return
	}
	user, err := getUser(emplID)
	if err != nil {
		log.Println("Failed to get users:", err)
		c.String(400, "")
		return
	}
	if deptID != strconv.Itoa(user.DeptID) {
		c.JSON(200, gin.H{"status": 0, "message": "Employee does not belong this department."})
		return
	}
	status := c.PostForm("status")
	if _, err := db.Exec(
		"UPDATE record SET user_id = ?, dept_id = ?, date = ?, type = ?, duration = ?, status = ?, describe = ? WHERE id = ?",
		emplID, deptID, date, Type, duration, status, describe, id); err != nil {
		log.Println("Failed to update record:", err)
		c.String(500, "")
		return
	}
	c.JSON(200, gin.H{"status": 1})
}

func verifyRecord(c *gin.Context) {
	if !verifyResponse("verify", c.ClientIP(), c.PostForm("g-recaptcha-response")) {
		c.String(403, "reCAPTCHA challenge failed")
		return
	}
	id := c.Param("id")
	records, _, err := getRecords(&idOptions{Record: id}, nil)
	if err != nil {
		log.Println("Failed to get record:", err)
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
		log.Println("Failed to get status:", err)
		c.String(400, "")
		return
	}
	if status != 1 && status != 2 {
		log.Println("Unknow status.")
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

	var user employee
	switch userID := sessions.Default(c).Get("userID"); userID {
	case "0":
		user = employee{ID: 0}
	default:
		user, err = getUser(userID)
		if err != nil {
			log.Println("Failed to get user:", err)
			c.String(500, "")
			return
		}
	}
	ip, _, _ := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr))
	if _, err := db.Exec("UPDATE record SET status = ?, verifiedby = ? WHERE id = ?",
		status, fmt.Sprintf("%d-%s-%s", user.ID, ip, c.ClientIP()), id); err != nil {
		log.Println("Failed to verify record:", err)
		c.String(500, "")
		return
	}
	localize := localize(c)
	var result string
	if status == 1 {
		result = localize["Verified"]
	} else {
		result = localize["Rejected"]
	}
	c.JSON(200, gin.H{"status": 1})
	notify(&idOptions{User: records[0].UserID},
		fmt.Sprintf(localize["VerifyRecordSubscribe"], user.Realname, result), localize)
	notify(&idOptions{Departments: []string{strconv.Itoa(records[0].DeptID)}},
		fmt.Sprintf(localize["VerifyRecordAdminSubscribe"], user.Realname, result), localize)
}

func deleteRecord(c *gin.Context) {
	id := c.Param("id")
	records, _, err := getRecords(&idOptions{Record: id}, nil)
	if err != nil {
		log.Println("Failed to get record:", err)
		c.String(400, "")
		return
	}
	var user employee
	switch userID := sessions.Default(c).Get("userID"); userID {
	case "0":
		user = employee{ID: 0}
	default:
		user, err = getUser(userID)
		if err != nil {
			log.Println("Failed to get user:", err)
			c.String(500, "")
			return
		}
	}
	if user.ID != 0 && records[0].Status != 0 {
		c.JSON(200, gin.H{"status": 0, "message": "You can only delete record which is not verified."})
		return
	}

	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	if _, err := db.Exec("DELETE FROM record WHERE id = ?", id); err != nil {
		log.Println("Failed to delete record:", err)
		c.String(500, "")
		return
	}
	c.JSON(200, gin.H{"status": 1})
	if user.ID != 0 {
		localize := localize(c)
		notify(&idOptions{Departments: []string{strconv.Itoa(user.DeptID)}},
			fmt.Sprintf(localize["DeleteRecordSubscribe"], user.Realname), localize)
	}
}
