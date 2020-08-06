package olms

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func get(c *gin.Context) {
	var user empl
	var total int
	var err error
	session := sessions.Default(c)
	switch userID := session.Get("userID"); userID {
	case "0":
		db, err := getDB()
		if err != nil {
			log.Printf("Failed to connect to database: %v", err)
			c.String(503, "")
			return
		}
		defer db.Close()
		var permission []byte
		if err := db.QueryRow("SELECT group_concat(id) FROM department").Scan(&permission); err != nil {
			log.Printf("Failed to get admin permission: %v", err)
			c.String(500, "")
			return
		}
		user = empl{ID: 0, Role: true, Permission: string(permission)}
	default:
		users, _, err := getEmpls(session.Get("userID"), nil, nil, nil)
		if err != nil {
			log.Printf("Failed to get user: %v", err)
			c.String(500, "")
			return
		}
		user = users[0]
	}
	query := c.PostForm("query")
	id := c.PostForm("id")
	userID := c.PostForm("user_id")
	deptID := c.PostForm("dept_id")
	period := c.PostForm("period")
	year := c.PostForm("year")
	month := c.PostForm("month")
	Type := c.PostForm("type")
	status := c.PostForm("status")
	role := c.PostForm("role")
	page := c.PostForm("page")

	switch c.PostForm("mode") {
	case "":
		if user.ID == 0 {
			log.Println("Super Administrator has no personal record.")
			c.String(400, "")
			return
		}
		switch query {
		case "records", "":
			var records []record
			if id != "" {
				records, total, err = getRecords(id, nil, nil, year, month, Type, status, page)
				if err != nil {
					log.Println(err)
					c.String(500, "")
					return
				}
				if records[0].UserID != user.ID {
					c.String(403, "")
					return
				}
				c.JSON(200, gin.H{"record": records[0]})
				return
			}
			records, total, err = getRecords(nil, user.ID, nil, year, month, Type, status, page)
			if err != nil {
				log.Println(err)
				c.String(500, "")
				return
			}
			c.JSON(200, gin.H{"total": total, "rows": records})
		case "stats":
			stats, total, err := getStats(user.ID, nil, period, year, month, page)
			if err != nil {
				log.Println(err)
				c.String(500, "")
				return
			}
			c.JSON(200, gin.H{"total": total, "rows": stats})
		default:
			c.String(400, "Unknown query")
		}
	case "admin":
		if !user.Role || user.ID != 0 {
			c.String(403, "")
			return
		}
		switch query {
		case "records", "":
			var records []record
			if id != "" {
				records, _, err = getRecords(id, nil, nil, year, month, Type, status, page)
				if err != nil {
					log.Println(err)
					c.String(500, "")
					return
				}
				for _, i := range strings.Split(user.Permission, ",") {
					if strconv.Itoa(records[0].DeptID) == i {
						c.JSON(200, gin.H{"record": records[0]})
						return
					}
				}
				c.String(403, "")
				return
			} else if userID != "" {
				if checkPermission(c, "", userID) {
					records, total, err = getRecords(nil, userID, nil, year, month, Type, status, page)
					if err != nil {
						log.Println(err)
						c.String(500, "")
						return
					}
				} else {
					c.String(403, "")
					return
				}
			} else if deptID != "" {
				if checkPermission(c, deptID) {
					records, total, err = getRecords(nil, nil, []string{deptID}, year, month, Type, status, page)
					if err != nil {
						log.Println(err)
						c.String(500, "")
						return
					}
				} else {
					c.String(403, "")
					return
				}
			} else {
				records, total, err = getRecords(nil, nil, strings.Split(user.Permission, ","), year, month, Type, status, page)
				if err != nil {
					log.Println(err)
					c.String(500, "")
					return
				}
			}
			c.JSON(200, gin.H{"total": total, "rows": records})
		case "stats":
			var stats []stat
			if userID != "" {
				if checkPermission(c, "", userID) {
					stats, total, err = getStats(userID, nil, period, year, month, page)
					if err != nil {
						log.Println(err)
						c.String(500, "")
						return
					}
				} else {
					c.String(403, "")
					return
				}
			} else if deptID != "" {
				if checkPermission(c, deptID) {
					stats, total, err = getStats(nil, []string{deptID}, period, year, month, page)
					if err != nil {
						log.Println(err)
						c.String(500, "")
						return
					}
				} else {
					c.String(403, "")
					return
				}
			} else {
				stats, total, err = getStats(nil, strings.Split(user.Permission, ","), period, year, month, page)
				if err != nil {
					log.Println(err)
					c.String(500, "")
					return
				}
			}
			c.JSON(200, gin.H{"total": total, "rows": stats})
		case "empls":
			var empls []empl
			if deptID != "" {
				if checkPermission(c, deptID) {
					empls, total, err = getEmpls(nil, []string{deptID}, role, page)
					if err != nil {
						log.Println(err)
						c.String(500, "")
						return
					}
				} else {
					c.String(403, "")
					return
				}
			} else {
				empls, total, err = getEmpls(nil, strings.Split(user.Permission, ","), role, page)
				if err != nil {
					log.Println(err)
					c.String(500, "")
					return
				}
			}
			for _, i := range empls {
				i.Role = false
				i.Permission = ""
			}
			c.JSON(200, gin.H{"total": total, "rows": empls})
		case "depts":
			depts, err := getDepts(strings.Split(user.Permission, ","))
			if err != nil {
				log.Println(err)
				c.String(500, "")
				return
			}
			c.JSON(200, gin.H{"rows": depts})
		default:
			c.String(400, "Unknown query")
		}
	case "super":
		if user.ID != 0 {
			c.String(403, "")
			return
		}
		switch query {
		case "empls", "":
			var empls []empl
			if deptID != "" {
				empls, total, err = getEmpls(nil, []string{deptID}, role, page)
				if err != nil {
					log.Println(err)
					c.String(500, "")
					return
				}
			} else {
				empls, total, err = getEmpls(nil, strings.Split(user.Permission, ","), role, page)
				if err != nil {
					log.Println(err)
					c.String(500, "")
					return
				}
			}
			c.JSON(200, gin.H{"total": total, "rows": empls})
		default:
			c.String(400, "Unknown query")
		}
	default:
		c.String(400, "Unknown query")
	}
}

func exportCSV(c *gin.Context) {
	var user empl
	var prefix string
	var results []map[string]interface{}
	var err error
	session := sessions.Default(c)
	switch userID := session.Get("userID"); userID {
	case "0":
		db, err := getDB()
		if err != nil {
			log.Printf("Failed to connect to database: %v", err)
			c.String(503, "")
			return
		}
		defer db.Close()
		var permission []byte
		if err := db.QueryRow("SELECT group_concat(id) FROM department").Scan(&permission); err != nil {
			log.Printf("Failed to get admin permission: %v", err)
			c.String(500, "")
			return
		}
		user = empl{ID: 0, Role: true, Permission: string(permission)}
	default:
		users, _, err := getEmpls(session.Get("userID"), nil, nil, nil)
		if err != nil {
			log.Printf("Failed to get user: %v", err)
			c.String(500, "")
			return
		}
		user = users[0]
	}
	query := c.PostForm("query")
	userID := c.PostForm("user_id")
	deptID := c.PostForm("dept_id")
	period := c.PostForm("period")
	year := c.PostForm("year")
	month := c.PostForm("month")
	Type := c.PostForm("type")
	status := c.PostForm("status")

	switch c.PostForm("mode") {
	case "":
		if user.ID == 0 {
			log.Println("Super Administrator has no personal record.")
			c.String(400, "")
			return
		}
		switch query {
		case "records", "":
			records, _, err := getRecords(nil, user.ID, nil, year, month, Type, status, nil)
			if err != nil {
				log.Println(err)
				c.String(500, "")
				return
			}
			for _, i := range records {
				results = append(results, i.format())
			}
			sendCSV(c,
				fmt.Sprintf("EmplRecords-%s%s%s.csv", user.Realname, year, month),
				[]string{"Date", "Type", "Duration", "Describe", "Created", "Status"},
				results)
		case "stats":
			stats, _, err := getStats(user.ID, nil, period, year, month, nil)
			if err != nil {
				log.Println(err)
				c.String(500, "")
				return
			}
			for _, i := range stats {
				results = append(results, i.format())
			}
			sendCSV(c,
				fmt.Sprintf("EmplStats-%s%s%s.csv", user.Realname, year, month),
				[]string{"Period", "DeptName", "Name", "Overtime", "Leave", "Summary"},
				results)
		default:
			c.String(400, "Unknown query")
		}
	case "admin":
		if !user.Role || user.ID != 0 {
			c.String(403, "")
			return
		}
		switch query {
		case "records", "":
			var records []record
			if userID != "" {
				if checkPermission(c, "", userID) {
					records, _, err = getRecords(nil, userID, nil, year, month, Type, status, nil)
					if err != nil {
						log.Println(err)
						c.String(500, "")
						return
					}
					if len(records) == 0 {
						c.String(200, "No result.")
						return
					}
					prefix = records[0].Name
				} else {
					c.String(403, "")
					return
				}
			} else if deptID != "" {
				if checkPermission(c, deptID) {
					records, _, err = getRecords(nil, nil, []string{deptID}, year, month, Type, status, nil)
					if err != nil {
						log.Println(err)
						c.String(500, "")
						return
					}
					if len(records) == 0 {
						c.String(200, "No result.")
						return
					}
					prefix = records[0].DeptName
				} else {
					c.String(403, "")
					return
				}
			} else {
				records, _, err = getRecords(nil, nil, strings.Split(user.Permission, ","), year, month, Type, status, nil)
				if err != nil {
					log.Println(err)
					c.String(500, "")
					return
				}
			}
			for _, i := range records {
				results = append(results, i.format())
			}
			sendCSV(c,
				fmt.Sprintf("DeptRecords%s%s%s.csv", prefix, year, month),
				[]string{"DeptName", "Name", "Date", "Type", "Duration", "Describe", "Created", "Status"},
				results)
		case "stats":
			var stats []stat
			if userID != "" {
				if checkPermission(c, "", userID) {
					stats, _, err = getStats(userID, nil, period, year, month, nil)
					if err != nil {
						log.Println(err)
						c.String(500, "")
						return
					}
					if len(stats) == 0 {
						c.String(200, "No result.")
						return
					}
					prefix = stats[0].Name
				} else {
					c.String(403, "")
					return
				}
			} else if deptID != "" {
				if checkPermission(c, deptID) {
					stats, _, err = getStats(nil, []string{deptID}, period, year, month, nil)
					if err != nil {
						log.Println(err)
						c.String(500, "")
						return
					}
					if len(stats) == 0 {
						c.String(200, "No result.")
						return
					}
					prefix = stats[0].DeptName
				} else {
					c.String(403, "")
					return
				}
			} else {
				stats, _, err = getStats(nil, strings.Split(user.Permission, ","), period, year, month, nil)
				if err != nil {
					log.Println(err)
					c.String(500, "")
					return
				}
			}
			for _, i := range stats {
				results = append(results, i.format())
			}
			sendCSV(c,
				fmt.Sprintf("DeptStats%s%s%s.csv", prefix, year, month),
				[]string{"Period", "DeptName", "Name", "Overtime", "Leave", "Summary"},
				results)
		default:
			c.String(400, "Unknown query")
		}
	default:
		c.String(400, "Unknown query")
	}
}
