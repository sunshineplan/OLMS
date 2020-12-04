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
	Record      interface{}
	User        interface{}
	Departments []string
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

	var user employee
	var err error
	switch userID := sessions.Default(c).Get("userID"); userID {
	case "0":
		db, err := getDB()
		if err != nil {
			log.Println("Failed to connect to database:", err)
			c.String(503, "")
			return
		}
		defer db.Close()
		var permission []byte
		if err := db.QueryRow("SELECT group_concat(id) FROM department").Scan(&permission); err != nil {
			log.Println("Failed to get admin permission:", err)
			c.String(500, "")
			return
		}
		user = employee{ID: 0, Role: true, Permission: string(permission)}
	default:
		user, err = getUser(userID)
		if err != nil {
			log.Println("Failed to get user:", err)
			c.String(500, "")
			return
		}
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
	switch obj["mode"] {
	case nil:
		if user.ID == 0 {
			log.Print("Super Administrator has no personal record.")
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
				records, _, err = getRecords(&idOptions{Record: id}, nil)
				if err != nil {
					log.Print(err)
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
			records, total, err = getRecords(&idOptions{User: user.ID}, &options)
			if err != nil {
				log.Print(err)
				c.String(500, "")
				return
			}
			c.JSON(200, gin.H{"total": total, "rows": records})
		case "stats":
			if options.Page == nil {
				options.Page = 1.0
			}
			stats, total, err := getStatistics(&idOptions{User: user.ID}, &options)
			if err != nil {
				log.Print(err)
				c.String(500, "")
				return
			}
			c.JSON(200, gin.H{"total": total, "rows": stats})
		case "years":
			years, err := getYears(&idOptions{User: user.ID})
			if err != nil {
				log.Print(err)
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
				records, _, err = getRecords(&idOptions{Record: id}, nil)
				if err != nil {
					log.Print(err)
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
				if checkPermission(c, &idOptions{User: userID}) {
					records, total, err = getRecords(&idOptions{User: userID}, &options)
					if err != nil {
						log.Print(err)
						c.String(500, "")
						return
					}
				} else {
					c.String(403, "")
					return
				}
			} else if deptID != nil {
				if checkPermission(c, &idOptions{Departments: []string{fmt.Sprintf("%v", deptID)}}) {
					records, total, err = getRecords(&idOptions{Departments: []string{fmt.Sprintf("%v", deptID)}}, &options)
					if err != nil {
						log.Print(err)
						c.String(500, "")
						return
					}
				} else {
					c.String(403, "")
					return
				}
			} else {
				records, total, err = getRecords(&idOptions{Departments: strings.Split(user.Permission, ",")}, &options)
				if err != nil {
					log.Print(err)
					c.String(500, "")
					return
				}
			}
			c.JSON(200, gin.H{"total": total, "rows": records})
		case "stats":
			if options.Page == nil {
				options.Page = 1.0
			}
			var statistics []statistic
			if userID != nil {
				if checkPermission(c, &idOptions{User: userID}) {
					statistics, total, err = getStatistics(&idOptions{User: userID}, &options)
					if err != nil {
						log.Print(err)
						c.String(500, "")
						return
					}
				} else {
					c.String(403, "")
					return
				}
			} else if deptID != nil {
				if checkPermission(c, &idOptions{Departments: []string{fmt.Sprintf("%v", deptID)}}) {
					statistics, total, err = getStatistics(&idOptions{Departments: []string{fmt.Sprintf("%v", deptID)}}, &options)
					if err != nil {
						log.Print(err)
						c.String(500, "")
						return
					}
				} else {
					c.String(403, "")
					return
				}
			} else {
				statistics, total, err = getStatistics(&idOptions{Departments: strings.Split(user.Permission, ",")}, &options)
				if err != nil {
					log.Print(err)
					c.String(500, "")
					return
				}
			}
			c.JSON(200, gin.H{"total": total, "rows": statistics})
		case "years":
			var years []string
			if userID != nil {
				if checkPermission(c, &idOptions{User: userID}) {
					years, err = getYears(&idOptions{User: userID})
					if err != nil {
						log.Print(err)
						c.String(500, "")
						return
					}
				} else {
					c.String(403, "")
					return
				}
			} else if deptID != nil {
				if checkPermission(c, &idOptions{Departments: []string{fmt.Sprintf("%v", deptID)}}) {
					years, err = getYears(&idOptions{Departments: []string{fmt.Sprintf("%v", deptID)}})
					if err != nil {
						log.Print(err)
						c.String(500, "")
						return
					}
				} else {
					c.String(403, "")
					return
				}
			} else {
				years, err = getYears(&idOptions{Departments: strings.Split(user.Permission, ",")})
				if err != nil {
					log.Print(err)
					c.String(500, "")
					return
				}
			}
			c.JSON(200, gin.H{"rows": years})
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

	var user employee
	var err error
	switch userID := sessions.Default(c).Get("userID"); userID {
	case "0":
		db, err := getDB()
		if err != nil {
			log.Println("Failed to connect to database:", err)
			c.String(503, "")
			return
		}
		defer db.Close()
		var permission []byte
		if err := db.QueryRow("SELECT group_concat(id) FROM department").Scan(&permission); err != nil {
			log.Println("Failed to get admin permission:", err)
			c.String(500, "")
			return
		}
		user = employee{ID: 0, Role: true, Permission: string(permission)}
	default:
		user, err = getUser(userID)
		if err != nil {
			log.Println("Failed to get user:", err)
			c.String(500, "")
			return
		}
	}

	localize := localize(c)

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
	switch obj["mode"] {
	case nil:
		if user.ID == 0 {
			log.Print("Super Administrator has no personal record.")
			c.String(400, "")
			return
		}
		switch query {
		case "records", nil:
			records, _, err := getRecords(&idOptions{User: user.ID}, &options)
			if err != nil {
				log.Print(err)
				c.String(500, "")
				return
			}
			for _, i := range records {
				results = append(results, i.format(localize))
			}
			sendCSV(c,
				fmt.Sprintf("%s-%s%v%v.csv", localize["EmplRecords"], user.Realname, options.Year, options.Month),
				[]string{
					localize["Date"],
					localize["Type"],
					localize["Duration"],
					localize["Describe"],
					localize["Created"],
					localize["Status"]},
				results)
		case "stats":
			stats, _, err := getStatistics(&idOptions{User: user.ID}, &options)
			if err != nil {
				log.Print(err)
				c.String(500, "")
				return
			}
			for _, i := range stats {
				results = append(results, i.format(localize))
			}
			sendCSV(c,
				fmt.Sprintf("%s-%s%v%v.csv", localize["EmplStats"], user.Realname, options.Year, options.Month),
				[]string{
					localize["Period"],
					localize["DeptName"],
					localize["Name"],
					localize["Overtime"],
					localize["Leave"],
					localize["Summary"]},
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
				if checkPermission(c, &idOptions{User: userID}) {
					records, _, err = getRecords(&idOptions{User: userID}, &options)
					if err != nil {
						log.Print(err)
						c.String(500, "")
						return
					}
					if len(records) == 0 {
						c.String(404, "No result.")
						return
					}
					prefix = records[0].Realname
				} else {
					c.String(403, "")
					return
				}
			} else if deptID != nil {
				if checkPermission(c, &idOptions{Departments: []string{fmt.Sprintf("%v", deptID)}}) {
					records, _, err = getRecords(&idOptions{Departments: []string{fmt.Sprintf("%v", deptID)}}, &options)
					if err != nil {
						log.Print(err)
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
				records, _, err = getRecords(&idOptions{Departments: strings.Split(user.Permission, ",")}, &options)
				if err != nil {
					log.Print(err)
					c.String(500, "")
					return
				}
			}
			for _, i := range records {
				results = append(results, i.format(localize))
			}
			sendCSV(c,
				fmt.Sprintf("%s%s%v%v.csv", localize["DeptRecords"], prefix, options.Year, options.Month),
				[]string{
					localize["DeptName"],
					localize["Name"],
					localize["Date"],
					localize["Type"],
					localize["Duration"],
					localize["Describe"],
					localize["Created"],
					localize["Status"]},
				results)
		case "stats":
			var statistics []statistic
			if userID != nil {
				if checkPermission(c, &idOptions{User: userID}) {
					statistics, _, err = getStatistics(&idOptions{User: userID}, &options)
					if err != nil {
						log.Print(err)
						c.String(500, "")
						return
					}
					if len(statistics) == 0 {
						c.String(404, "No result.")
						return
					}
					prefix = statistics[0].Realname
				} else {
					c.String(403, "")
					return
				}
			} else if deptID != nil {
				if checkPermission(c, &idOptions{Departments: []string{fmt.Sprintf("%v", deptID)}}) {
					statistics, _, err = getStatistics(&idOptions{Departments: []string{fmt.Sprintf("%v", deptID)}}, &options)
					if err != nil {
						log.Print(err)
						c.String(500, "")
						return
					}
					if len(statistics) == 0 {
						c.String(404, "No result.")
						return
					}
					prefix = statistics[0].DeptName
				} else {
					c.String(403, "")
					return
				}
			} else {
				statistics, _, err = getStatistics(&idOptions{Departments: strings.Split(user.Permission, ",")}, &options)
				if err != nil {
					log.Print(err)
					c.String(500, "")
					return
				}
			}
			for _, i := range statistics {
				results = append(results, i.format(localize))
			}
			sendCSV(c,
				fmt.Sprintf("%s%s%v%v.csv", localize["DeptStats"], prefix, options.Year, options.Month),
				[]string{
					localize["Period"],
					localize["DeptName"],
					localize["Name"],
					localize["Overtime"],
					localize["Leave"],
					localize["Summary"]},
				results)
		default:
			c.String(400, "Unknown query")
		}
	default:
		c.String(400, "Unknown query")
	}
}
