package olms

import (
	"log"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func get(c *gin.Context) {
	var total int
	var err error
	session := sessions.Default(c)
	users, _, err := getEmpls(session.Get("userID"), nil, nil, nil)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		c.String(500, "")
		return
	}
	user := users[0]
	query := c.PostForm("query")
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
			records, total, err := getRecords(user.ID, nil, year, month, Type, status, page)
			if err != nil {
				log.Println(err)
				c.String(500, "")
				return
			}
			c.JSON(200, gin.H{"total": total, "records": records})
		case "stats":
			stats, total, err := getStats(user.ID, nil, period, year, month, page)
			if err != nil {
				log.Println(err)
				c.String(500, "")
				return
			}
			c.JSON(200, gin.H{"total": total, "statistics": stats})
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
					records, total, err = getRecords(userID, nil, year, month, Type, status, page)
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
					records, total, err = getRecords(nil, []string{deptID}, year, month, Type, status, page)
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
				records, total, err = getRecords(nil, strings.Split(user.Permission, ","), year, month, Type, status, page)
				if err != nil {
					log.Println(err)
					c.String(500, "")
					return
				}
			}
			c.JSON(200, gin.H{"total": total, "records": records})
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
			c.JSON(200, gin.H{"total": total, "statistics": stats})
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
			c.JSON(200, gin.H{"total": total, "employees": empls})
		case "depts":
			depts, err := getDepts(strings.Split(user.Permission, ","))
			if err != nil {
				log.Println(err)
				c.String(500, "")
				return
			}
			c.JSON(200, gin.H{"departments": depts})
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
			c.JSON(200, gin.H{"total": total, "employees": empls})
		default:
			c.String(400, "Unknown query")
		}
	default:
		c.String(400, "Unknown query")
	}
}
