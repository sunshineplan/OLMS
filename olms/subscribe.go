package olms

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/utils/workers"
)

func isEmailValid(email string) bool {
	if len(email) < 3 && len(email) > 254 {
		return false
	}
	//https://www.w3.org/TR/2016/REC-html51-20161101/sec-forms.html#email-state-typeemail
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(email)
}

func subscribe(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	userID := sessions.Default(c).Get("userID")
	switch c.PostForm("subscribe") {
	case "0":
		if _, err := db.Exec("UPDATE user SET subscribe = 0 WHERE id = ?", userID); err != nil {
			log.Printf("Failed to update user: %v", err)
			c.String(500, "")
			return
		}
	case "1":
		if email := c.PostForm("email"); isEmailValid(email) {
			if _, err := db.Exec(
				"UPDATE user SET subscribe = 1, email = ? WHERE id = ?", email, userID); err != nil {
				log.Printf("Failed to update user: %v", err)
				c.String(500, "")
				return
			}
		} else {
			c.JSON(200, gin.H{"status": 0})
			return
		}
	default:
		c.String(400, "")
		return
	}
	c.JSON(200, gin.H{"status": 1})
}

func notify(id *idOptions, message string, localize translate) {
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	var emails []interface{}
	var title string
	if id.UserID != nil {
		var subscribe bool
		var realname, email string
		if err := db.QueryRow(
			"SELECT subscribe, realname, email FROM user WHERE id = ?",
			id.UserID).Scan(&subscribe, &realname, &email); err != nil {
			log.Printf("Failed to get user subscribe: %v", err)
		}
		if subscribe {
			title = fmt.Sprintf(localize["Dear"], realname)
			emails = append(emails, email)
		} else {
			return
		}
	} else if len(id.DeptIDs) == 1 {
		var deptName string
		if err := db.QueryRow("SELECT dept_name FROM department WHERE id = ?",
			id.DeptIDs[0]).Scan(&deptName); err != nil {
			log.Printf("Failed to get department: %v", err)
			return
		}
		rows, err := db.Query(
			"SELECT email FROM user LEFT JOIN permission p ON id = user_id WHERE subscribe = 1 AND (p.dept_id = ? OR id = 0)",
			id.DeptIDs[0])
		if err != nil {
			log.Printf("Failed to get employees: %v", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var email string
			if err = rows.Scan(&email); err != nil {
				log.Printf("Failed to scan email: %v", err)
				return
			}
			emails = append(emails, email)
		}
		title = fmt.Sprintf(localize["DearAdmin"], deptName)
	} else {
		return
	}
	workers.New(5).Run(emails, func(c chan bool, _ int, email interface{}) {
		defer func() { <-c }()
		subscribe := MailSetting
		subscribe.To = []string{email.(string)}
		if err := subscribe.Send(
			fmt.Sprintf("%s-%s", localize["NotificationSubject"], time.Now().Format("20060102 15:04")),
			fmt.Sprintf("%s\n\n    %s", title, message),
		); err != nil {
			log.Printf("Failed to send mail: %v", err)
		}
	})
}
