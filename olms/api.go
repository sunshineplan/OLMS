package olms

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type idOptions struct {
	RecordID interface{}
	UserID   interface{}
	DeptIDs  []string
}

type searchOptions struct {
	Period   interface{}
	Year     interface{}
	Month    interface{}
	Type     interface{}
	Role     interface{}
	Status   interface{}
	Describe interface{}
	Page     interface{}
	Sort     interface{}
	Order    interface{}
}

func get(c *gin.Context) {
	var obj map[string]interface{}
	if err := c.BindJSON(&obj); err != nil {
		c.String(400, "")
		return
	}

	if !verifyResponse("get", c.ClientIP(), obj["g-recaptcha-response"]) {
		c.String(403, "reCAPTCHA challenge failed")
		return
	}

	var user empl
	switch userID := sessions.Default(c).Get("userID"); userID {
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
		users, _, err := getEmpls(&idOptions{UserID: userID}, nil)
		if err != nil {
			log.Printf("Failed to get user: %v", err)
			c.String(500, "")
			return
		}
		user = users[0]
	}

	var options searchOptions
	query := obj["query"]
	id := obj["id"]
	userID := obj["empl"]
	deptID := obj["dept"]
	options.Period = obj["period"]
	options.Year = obj["year"]
	options.Month = obj["month"]
	options.Type = obj["type"]
	options.Status = obj["status"]
	options.Describe = obj["describe"]
	options.Role = obj["role"]
	options.Page = obj["page"]
	options.Sort = obj["sort"]
	options.Order = obj["order"]

	var total int
	var err error
	switch obj["mode"] {
	case nil:
		if user.ID == 0 {
			log.Println("Super Administrator has no personal record.")
			c.String(400, "")
			return
		}
		switch query {
		case "records", nil:
			if options.Page == nil {
				options.Page = 1.0
			}
			var records []record
			if id != nil {
				records, _, err = getRecords(&idOptions{RecordID: id}, nil)
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
			records, total, err = getRecords(&idOptions{UserID: user.ID}, &options)
			if err != nil {
				log.Println(err)
				c.String(500, "")
				return
			}
			c.JSON(200, gin.H{"total": total, "rows": records})
		case "stats":
			if options.Page == nil {
				options.Page = 1.0
			}
			stats, total, err := getStats(&idOptions{UserID: user.ID}, &options)
			if err != nil {
				log.Println(err)
				c.String(500, "")
				return
			}
			c.JSON(200, gin.H{"total": total, "rows": stats})
		case "years":
			years, err := getYears(&idOptions{UserID: user.ID})
			if err != nil {
				log.Println(err)
				c.String(500, "")
				return
			}
			c.JSON(200, gin.H{"rows": years})
		default:
			c.String(400, "Unknown query")
		}
	case "admin":
		if !user.Role && user.ID != 0 {
			c.String(403, "")
			return
		}
		switch query {
		case "records", nil:
			if options.Page == nil {
				options.Page = 1.0
			}
			var records []record
			if id != nil {
				records, _, err = getRecords(&idOptions{RecordID: id}, nil)
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
			} else if userID != nil {
				if checkPermission(c, nil, userID) {
					records, total, err = getRecords(&idOptions{UserID: userID}, &options)
					if err != nil {
						log.Println(err)
						c.String(500, "")
						return
					}
				} else {
					c.String(403, "")
					return
				}
			} else if deptID != nil {
				if checkPermission(c, deptID) {
					records, total, err = getRecords(&idOptions{DeptIDs: []string{fmt.Sprintf("%v", deptID)}}, &options)
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
				records, total, err = getRecords(&idOptions{DeptIDs: strings.Split(user.Permission, ",")}, &options)
				if err != nil {
					log.Println(err)
					c.String(500, "")
					return
				}
			}
			c.JSON(200, gin.H{"total": total, "rows": records})
		case "stats":
			if options.Page == nil {
				options.Page = 1.0
			}
			var stats []stat
			if userID != nil {
				if checkPermission(c, nil, userID) {
					stats, total, err = getStats(&idOptions{UserID: userID}, &options)
					if err != nil {
						log.Println(err)
						c.String(500, "")
						return
					}
				} else {
					c.String(403, "")
					return
				}
			} else if deptID != nil {
				if checkPermission(c, deptID) {
					stats, total, err = getStats(&idOptions{DeptIDs: []string{fmt.Sprintf("%v", deptID)}}, &options)
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
				stats, total, err = getStats(&idOptions{DeptIDs: strings.Split(user.Permission, ",")}, &options)
				if err != nil {
					log.Println(err)
					c.String(500, "")
					return
				}
			}
			c.JSON(200, gin.H{"total": total, "rows": stats})
		case "empls":
			var empls []empl
			if id != nil {
				empls, _, err = getEmpls(&idOptions{UserID: id}, nil)
				if err != nil {
					log.Println(err)
					c.String(500, "")
					return
				}
				for _, i := range strings.Split(user.Permission, ",") {
					if strconv.Itoa(empls[0].DeptID) == i {
						c.JSON(200, gin.H{"empl": empls[0]})
						return
					}
				}
				c.String(403, "")
				return
			} else if deptID != nil {
				if checkPermission(c, deptID) {
					empls, total, err = getEmpls(&idOptions{DeptIDs: []string{fmt.Sprintf("%v", deptID)}}, &options)
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
				empls, total, err = getEmpls(&idOptions{DeptIDs: strings.Split(user.Permission, ",")}, &options)
				if err != nil {
					log.Println(err)
					c.String(500, "")
					return
				}
			}
			for i := range empls {
				empls[i].Role = false
				empls[i].Permission = ""
			}
			c.JSON(200, gin.H{"total": total, "rows": empls})
		case "depts":
			var depts []dept
			if id != nil {
				depts, err = getDepts([]string{fmt.Sprintf("%v", id)})
				if err != nil {
					log.Println(err)
					c.String(500, "")
					return
				}
				for _, i := range strings.Split(user.Permission, ",") {
					if strconv.Itoa(depts[0].ID) == i {
						c.JSON(200, gin.H{"dept": depts[0]})
						return
					}
				}
				c.String(403, "")
				return
			}
			depts, err := getDepts(strings.Split(user.Permission, ","))
			if err != nil {
				log.Println(err)
				c.String(500, "")
				return
			}
			c.JSON(200, gin.H{"rows": depts})
		case "years":
			var years []string
			if userID != nil {
				if checkPermission(c, nil, userID) {
					years, err = getYears(&idOptions{UserID: userID})
					if err != nil {
						log.Println(err)
						c.String(500, "")
						return
					}
				} else {
					c.String(403, "")
					return
				}
			} else if deptID != nil {
				if checkPermission(c, deptID) {
					years, err = getYears(&idOptions{DeptIDs: []string{fmt.Sprintf("%v", deptID)}})
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
				years, err = getYears(&idOptions{DeptIDs: strings.Split(user.Permission, ",")})
				if err != nil {
					log.Println(err)
					c.String(500, "")
					return
				}
			}
			c.JSON(200, gin.H{"rows": years})
		case "info":
			depts, err := getDepts(strings.Split(user.Permission, ","))
			if err != nil {
				log.Println(err)
				c.String(500, "")
				return
			}
			empls, _, err := getEmpls(&idOptions{DeptIDs: strings.Split(user.Permission, ",")}, nil)
			if err != nil {
				log.Println(err)
				c.String(500, "")
				return
			}
			if user.ID != 0 {
				for i := range empls {
					empls[i].Role = false
					empls[i].Permission = ""
				}
			}
			years, err := getYears(&idOptions{DeptIDs: strings.Split(user.Permission, ",")})
			if err != nil {
				log.Println(err)
				c.String(500, "")
				return
			}
			c.JSON(200, gin.H{"depts": depts, "empls": empls, "years": years})
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
			if id != nil {
				empls, _, err = getEmpls(&idOptions{UserID: id}, nil)
				if err != nil {
					log.Println(err)
					c.String(500, "")
					return
				}
				c.JSON(200, gin.H{"empl": empls[0]})
				return
			} else if deptID != nil {
				empls, total, err = getEmpls(&idOptions{DeptIDs: []string{fmt.Sprintf("%v", deptID)}}, &options)
				if err != nil {
					log.Println(err)
					c.String(500, "")
					return
				}
			} else {
				empls, total, err = getEmpls(&idOptions{DeptIDs: strings.Split(user.Permission, ",")}, &options)
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
	var obj map[string]interface{}
	if err := c.BindJSON(&obj); err != nil {
		c.String(400, "")
		return
	}

	if !verifyResponse("export", c.ClientIP(), obj["g-recaptcha-response"]) {
		c.String(403, "reCAPTCHA challenge failed")
		return
	}

	var user empl
	switch userID := sessions.Default(c).Get("userID"); userID {
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
		users, _, err := getEmpls(&idOptions{UserID: userID}, nil)
		if err != nil {
			log.Printf("Failed to get user: %v", err)
			c.String(500, "")
			return
		}
		user = users[0]
	}

	var options searchOptions
	query := obj["query"]
	userID := obj["empl"]
	deptID := obj["dept"]
	options.Period = obj["period"]
	options.Year = obj["year"]
	options.Month = obj["month"]
	options.Type = obj["type"]
	options.Status = obj["status"]
	options.Describe = obj["describe"]

	var prefix string
	var results []map[string]interface{}
	var err error
	switch obj["mode"] {
	case nil:
		if user.ID == 0 {
			log.Println("Super Administrator has no personal record.")
			c.String(400, "")
			return
		}
		switch query {
		case "records", nil:
			records, _, err := getRecords(&idOptions{UserID: user.ID}, &options)
			if err != nil {
				log.Println(err)
				c.String(500, "")
				return
			}
			for _, i := range records {
				results = append(results, i.format())
			}
			sendCSV(c,
				fmt.Sprintf("EmplRecords-%s%v%v.csv", user.Realname, options.Year, options.Month),
				[]string{"Date", "Type", "Duration", "Describe", "Created", "Status"},
				results)
		case "stats":
			stats, _, err := getStats(&idOptions{UserID: user.ID}, &options)
			if err != nil {
				log.Println(err)
				c.String(500, "")
				return
			}
			for _, i := range stats {
				results = append(results, i.format())
			}
			sendCSV(c,
				fmt.Sprintf("EmplStats-%s%v%v.csv", user.Realname, options.Year, options.Month),
				[]string{"Period", "DeptName", "Name", "Overtime", "Leave", "Summary"},
				results)
		default:
			c.String(400, "Unknown query")
		}
	case "admin":
		if !user.Role && user.ID != 0 {
			c.String(403, "")
			return
		}
		switch query {
		case "records", nil:
			var records []record
			if userID != nil {
				if checkPermission(c, nil, userID) {
					records, _, err = getRecords(&idOptions{UserID: userID}, &options)
					if err != nil {
						log.Println(err)
						c.String(500, "")
						return
					}
					if len(records) == 0 {
						c.String(404, "No result.")
						return
					}
					prefix = records[0].Name
				} else {
					c.String(403, "")
					return
				}
			} else if deptID != nil {
				if checkPermission(c, deptID) {
					records, _, err = getRecords(&idOptions{DeptIDs: []string{fmt.Sprintf("%v", deptID)}}, &options)
					if err != nil {
						log.Println(err)
						c.String(500, "")
						return
					}
					if len(records) == 0 {
						c.String(404, "No result.")
						return
					}
					prefix = records[0].DeptName
				} else {
					c.String(403, "")
					return
				}
			} else {
				records, _, err = getRecords(&idOptions{DeptIDs: strings.Split(user.Permission, ",")}, &options)
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
				fmt.Sprintf("DeptRecords%s%v%v.csv", prefix, options.Year, options.Month),
				[]string{"DeptName", "Name", "Date", "Type", "Duration", "Describe", "Created", "Status"},
				results)
		case "stats":
			var stats []stat
			if userID != nil {
				if checkPermission(c, nil, userID) {
					stats, _, err = getStats(&idOptions{UserID: userID}, &options)
					if err != nil {
						log.Println(err)
						c.String(500, "")
						return
					}
					if len(stats) == 0 {
						c.String(404, "No result.")
						return
					}
					prefix = stats[0].Name
				} else {
					c.String(403, "")
					return
				}
			} else if deptID != nil {
				if checkPermission(c, deptID) {
					stats, _, err = getStats(&idOptions{DeptIDs: []string{fmt.Sprintf("%v", deptID)}}, &options)
					if err != nil {
						log.Println(err)
						c.String(500, "")
						return
					}
					if len(stats) == 0 {
						c.String(404, "No result.")
						return
					}
					prefix = stats[0].DeptName
				} else {
					c.String(403, "")
					return
				}
			} else {
				stats, _, err = getStats(&idOptions{DeptIDs: strings.Split(user.Permission, ",")}, &options)
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
				fmt.Sprintf("DeptStats%s%v%v.csv", prefix, options.Year, options.Month),
				[]string{"Period", "DeptName", "Name", "Overtime", "Leave", "Summary"},
				results)
		default:
			c.String(400, "Unknown query")
		}
	default:
		c.String(400, "Unknown query")
	}
}
