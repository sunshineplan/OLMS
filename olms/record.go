package olms

import (
	"database/sql"
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

func getRecords(db *sql.DB, id *idOptions, options *searchOptions) ([]record, int, error) {
	stmt := "SELECT %s FROM record JOIN employee ON user_id = employee.id WHERE"

	var args []interface{}
	var orderBy, limit string
	c := make(chan error, 1)
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

	if options.Year != "" {
		if options.Month == "" {
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
	if options.Describe != "" {
		stmt += " AND describe LIKE ?"
		args = append(args, fmt.Sprintf("%%%v%%", options.Describe))
	}
	if options.Page != 0 {
		limit = fmt.Sprintf(" LIMIT ?, ?")
		args = append(args, int(options.Page-1)*perPage, perPage)
	} else {
		c <- nil
	}
	if options.Sort != "" {
		orderBy = fmt.Sprintf(" ORDER BY %v %v", options.Sort, options.Order)
	} else {
		orderBy = " ORDER BY created DESC"
	}

	var total int
	go func() {
		c <- db.QueryRow(fmt.Sprintf(stmt, "count(*)"), args...).Scan(&total)
	}()

	rows, err := db.Query(fmt.Sprintf(stmt+orderBy+limit,
		"record.id, employee.dept_id, deptname, employee.id, realname, date, ABS(duration), type, status, describe, created"),
		args...)
	if err != nil {
		log.Println("Failed to get records:", err)
		return nil, 0, err
	}
	defer rows.Close()

	records := []record{}
	for rows.Next() {
		var r record
		if err := rows.Scan(
			&r.ID, &r.DeptID, &r.DeptName, &r.UserID, &r.Realname, &r.Date, &r.Duration, &r.Type, &r.Status, &r.Describe, &r.Created,
		); err != nil {
			log.Println("Failed to scan record:", err)
			return nil, 0, err
		}
		records = append(records, r)
	}
	if err := <-c; err != nil {
		log.Println("Failed to get total records:", err)
		return nil, 0, err
	}
	return records, total, nil
}

func addRecord(c *gin.Context) {
	var r struct {
		record
		Recaptcha string
	}
	if err := c.BindJSON(&r); err != nil {
		log.Println("Failed to get option:", err)
		c.String(400, "")
		return
	}

	if !verifyResponse("record", c.ClientIP(), r.Recaptcha) {
		c.String(403, "reCAPTCHAChallengeFailed")
		return
	}

	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	if r.Type {
		if r.Duration < 1 {
			c.JSON(200, gin.H{"status": 0, "message": "BadDuration"})
			return
		}
	} else {
		r.Duration = -r.Duration
		if r.Duration > -1 {
			c.JSON(200, gin.H{"status": 0, "message": "BadDuration"})
			return
		}
	}

	user, err := getUser(db, sessions.Default(c).Get("userID"))
	if err != nil {
		log.Println("Failed to get user:", err)
		c.String(500, "")
		return
	}

	ip, _, _ := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr))
	ip = ip + "-" + c.ClientIP()
	localize := localize(c)
	if r.UserID == 0 {
		if _, err := db.Exec(
			`INSERT INTO record (date, type, duration, describe, dept_id, user_id, createdby)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			r.Date, r.Type, r.Duration, r.Describe, user.DeptID, user.ID, fmt.Sprintf("%d-%s", user.ID, ip),
		); err != nil {
			log.Println("Failed to add record:", err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		notify(&idOptions{Departments: []string{strconv.Itoa(user.DeptID)}},
			fmt.Sprintf(localize["AddRecordSubscribe"], user.Realname), localize)
		return
	}
	if user.ID != 0 {
		if r.DeptID != 0 && !checkPermission(db, c, &idOptions{User: r.UserID, Departments: []string{strconv.Itoa(r.DeptID)}}) {
			c.String(403, "")
			return
		}
		if _, err := db.Exec(
			`INSERT INTO record (dept_id, user_id, date, type, duration, describe, status, createdby, verifiedby)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			r.DeptID, r.UserID, r.Date, r.Type, r.Duration, r.Describe, 1,
			fmt.Sprintf("%d-%s", user.ID, ip), fmt.Sprintf("%d-%s", user.ID, ip),
		); err != nil {
			log.Println("Failed to add record:", err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		notify(&idOptions{User: r.UserID},
			fmt.Sprintf(localize["AdminAddRecordSubscribe"], user.Realname), localize)
		notify(&idOptions{Departments: []string{strconv.Itoa(r.DeptID)}},
			fmt.Sprintf(localize["AdminAddRecordAdminSubscribe"], user.Realname), localize)
		return
	}
	if _, err := db.Exec(
		`INSERT INTO record (dept_id, user_id, date, type, duration, describe, status, createdby, verifiedby)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		r.DeptID, r.UserID, r.Date, r.Type, r.Duration, r.Describe, r.Status,
		fmt.Sprintf("%d-%s", 0, ip), fmt.Sprintf("%d-%s", 0, ip),
	); err != nil {
		log.Println("Failed to add record:", err)
		c.String(500, "")
		return
	}
	c.JSON(200, gin.H{"status": 1})
}

func editRecord(c *gin.Context) {
	var r struct {
		record
		Recaptcha string
	}
	if err := c.BindJSON(&r); err != nil {
		log.Println("Failed to get option:", err)
		c.String(400, "")
		return
	}

	if !verifyResponse("record", c.ClientIP(), r.Recaptcha) {
		c.String(403, "reCAPTCHAChallengeFailed")
		return
	}

	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	if r.Type {
		if r.Duration < 1 {
			c.JSON(200, gin.H{"status": 0, "message": "BadDuration"})
			return
		}
	} else {
		r.Duration = -r.Duration
		if r.Duration > -1 {
			c.JSON(200, gin.H{"status": 0, "message": "BadDuration"})
			return
		}
	}

	record, ok := checkRecord(db, c, r.ID, true)
	if !ok {
		c.String(403, "")
		return
	}
	userID := sessions.Default(c).Get("userID")
	if userID != 0 {
		if record.Status != 0 {
			c.JSON(200, gin.H{"status": 0, "message": "UpdateRecordNotVerified"})
			return
		}
		if _, err := db.Exec("UPDATE record SET date = ?, type = ?, duration = ?, describe = ? WHERE id = ?",
			r.Date, r.Type, r.Duration, r.Describe, r.ID); err != nil {
			log.Println("Failed to update record:", err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		localize := localize(c)
		notify(&idOptions{Departments: []string{strconv.Itoa(record.DeptID)}},
			fmt.Sprintf(localize["EditRecordSubscribe"], record.Realname), localize)
		return
	}
	if r.UserID == 0 || r.DeptID == 0 {
		log.Print("Missing param.")
		c.String(400, "")
		return
	}
	user, err := getUser(db, r.UserID)
	if err != nil {
		log.Println("Failed to get users:", err)
		c.String(400, "")
		return
	}
	if r.DeptID != user.DeptID {
		c.JSON(200, gin.H{"status": 0, "message": "EmployeeNotBelong"})
		return
	}
	if _, err := db.Exec(
		"UPDATE record SET user_id = ?, dept_id = ?, date = ?, type = ?, duration = ?, status = ?, describe = ? WHERE id = ?",
		r.UserID, r.DeptID, r.Date, r.Type, r.Duration, r.Status, r.Describe, r.ID); err != nil {
		log.Println("Failed to update record:", err)
		c.String(500, "")
		return
	}
	c.JSON(200, gin.H{"status": 1})
}

func verifyRecord(c *gin.Context) {
	var r struct {
		Status    bool
		Recaptcha string
	}
	if err := c.BindJSON(&r); err != nil {
		log.Println("Failed to get option:", err)
		c.String(400, "")
		return
	}

	if !verifyResponse("verify", c.ClientIP(), r.Recaptcha) {
		c.String(403, "reCAPTCHAChallengeFailed")
		return
	}

	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("Failed to get id:", err)
		c.String(400, "")
		return
	}
	record, ok := checkRecord(db, c, id, false)
	if !ok {
		c.String(403, "")
		return
	}
	if record.Status != 0 {
		log.Println("The record is already verified.")
		c.String(400, "")
		return
	}
	var status int
	if r.Status {
		status = 1
	} else {
		status = 2
	}

	user, err := getUser(db, sessions.Default(c).Get("userID"))
	if err != nil {
		log.Println("Failed to get user:", err)
		c.String(500, "")
		return
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
	if r.Status {
		result = localize["Verified"]
	} else {
		result = localize["Rejected"]
	}
	c.JSON(200, gin.H{"status": 1})
	notify(&idOptions{User: record.UserID},
		fmt.Sprintf(localize["VerifyRecordSubscribe"], user.Realname, result), localize)
	notify(&idOptions{Departments: []string{strconv.Itoa(record.DeptID)}},
		fmt.Sprintf(localize["VerifyRecordAdminSubscribe"], user.Realname, result), localize)
}

func deleteRecord(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("Failed to get id:", err)
		c.String(400, "")
		return
	}
	record, ok := checkRecord(db, c, id, true)
	if !ok {
		c.String(403, "")
		return
	}

	user, err := getUser(db, sessions.Default(c).Get("userID"))
	if err != nil {
		log.Println("Failed to get user:", err)
		c.String(500, "")
		return
	}
	if user.ID != 0 && record.Status != 0 {
		c.JSON(200, gin.H{"status": 0, "message": "DeleteRecordNotVerified"})
		return
	}

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
