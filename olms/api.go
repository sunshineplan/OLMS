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
	UserID    int
	DeptID    int
	Period    string
	Year      string
	Month     string
	Type      interface{}
	Role      interface{}
	Status    interface{}
	Describe  string
	Page      int
	Sort      string
	Order     string
	Recaptcha string
}

func info(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	userID := sessions.Default(c).Get("userID")
	user, err := getUser(db, userID)
	if err != nil {
		log.Println("Failed to get user:", err)
		c.String(500, "")
		return
	}
	if user.Role {
		var departments []department
		var employees []employee
		if user.ID == 0 {
			departments, err = getDepartments(db, nil, true)
			if err != nil {
				log.Println("Failed to get departments:", err)
				c.String(500, "")
				return
			}
			employees, err = getEmployees(db, nil, true)
			if err != nil {
				log.Println("Failed to get employees:", err)
				c.String(500, "")
				return
			}
		} else {
			departments, err = getDepartments(db, strings.Split(user.Permission, ","), false)
			if err != nil {
				log.Println("Failed to get departments:", err)
				c.String(500, "")
				return
			}
			employees, err = getEmployees(db, strings.Split(user.Permission, ","), false)
			if err != nil {
				log.Println("Failed to get employees:", err)
				c.String(500, "")
				return
			}
		}
		c.JSON(200, gin.H{"user": user, "departments": departments, "employees": employees})
	} else {
		c.JSON(200, gin.H{"user": user})
	}
}

func api(c *gin.Context, mode string, export bool) {
	var option searchOptions
	if err := c.BindJSON(&option); err != nil {
		c.String(400, "")
		return
	}

	if !verifyResponse(mode, c.ClientIP(), option.Recaptcha) {
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

	var ids idOptions
	if option.UserID != 0 {
		ids.User = option.UserID
	} else if option.DeptID != 0 {
		ids.Departments = []string{strconv.Itoa(option.DeptID)}
	} else {
		user, err := getUser(db, sessions.Default(c).Get("userID"))
		if err != nil {
			log.Println("Failed to get user:", err)
			c.String(500, "")
			return
		}
		ids.Departments = strings.Split(user.Permission, ",")
	}
	if checkPermission(db, c, &ids) {
		if option.Page == 0 {
			option.Page = 1
		} else if export {
			option.Page = 0
		}
		localize := localize(c)
		if mode == "records" {
			records, total, err := getRecords(db, &ids, &option)
			if err != nil {
				log.Println("Failed to get records:", err)
				c.String(500, "")
				return
			}
			if export {
				result := make([]map[string]interface{}, len(records))
				for i := range records {
					result[i] = records[i].format(localize)
				}
				sendCSV(c,
					fmt.Sprintf("%s%s%s.csv", localize["Records"], option.Year, option.Month),
					[]string{
						localize["DeptName"],
						localize["Name"],
						localize["Date"],
						localize["Type"],
						localize["Duration"],
						localize["Describe"],
						localize["Created"],
						localize["Status"]},
					result)
				return
			}
			c.JSON(200, gin.H{"rows": records, "total": total})
			return
		}
		if mode == "statistics" {
			statistics, total, err := getStatistics(db, &ids, &option)
			if err != nil {
				log.Println("Failed to get statistics:", err)
				c.String(500, "")
				return
			}
			if export {
				result := make([]map[string]interface{}, len(statistics))
				for i := range statistics {
					result[i] = statistics[i].format(localize)
				}
				sendCSV(c,
					fmt.Sprintf("%s%s%s.csv", localize["Statistics"], option.Year, option.Month),
					[]string{
						localize["Period"],
						localize["DeptName"],
						localize["Name"],
						localize["Overtime"],
						localize["Leave"],
						localize["Summary"]},
					result)
				return
			}
			c.JSON(200, gin.H{"rows": statistics, "total": total})
			return
		}
	}
	c.String(403, "")
}

func records(c *gin.Context) {
	api(c, "records", false)
}

func statistics(c *gin.Context) {
	api(c, "statistics", false)
}

func exportRecords(c *gin.Context) {
	api(c, "records", true)
}

func exportStatistics(c *gin.Context) {
	api(c, "statistics", true)
}
