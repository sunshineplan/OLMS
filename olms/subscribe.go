package olms

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/utils/mail"
	"github.com/sunshineplan/utils/workers"
)

func isEmailValid(email string) bool {
	if len(email) < 3 && len(email) > 254 {
		return false
	}
	//https://www.w3.org/TR/2016/REC-html51-20161101/sec-forms.html#email-state-typeemail
	emailRegex := regexp.MustCompile(
		"^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(email)
}

func getSubscribe(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	var subscribe bool
	var email string
	if err := db.QueryRow(
		"SELECT subscribe, email FROM user WHERE id = ?",
		sessions.Default(c).Get("userID")).Scan(&subscribe, &email); err != nil {
		log.Println("Failed to get subscribe:", err)
		c.String(500, "")
		return
	}
	c.JSON(200, gin.H{"subscribe": subscribe, "Email": email})
}

func subscribe(c *gin.Context) {
	var subscribe struct {
		Subscribe bool
		Email     string
	}
	if err := c.BindJSON(&subscribe); err != nil {
		log.Println("Failed to get option:", err)
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

	userID := sessions.Default(c).Get("userID")
	if subscribe.Subscribe {
		if isEmailValid(subscribe.Email) {
			if _, err := db.Exec(
				"UPDATE user SET subscribe = 1, email = ? WHERE id = ?", subscribe.Email, userID); err != nil {
				log.Println("Failed to update user:", err)
				c.String(500, "")
				return
			}
		} else {
			c.JSON(200, gin.H{"status": 0})
			return
		}
	} else {
		if _, err := db.Exec("UPDATE user SET subscribe = 0 WHERE id = ?", userID); err != nil {
			log.Println("Failed to update user:", err)
			c.String(500, "")
			return
		}
	}
	c.JSON(200, gin.H{"status": 1})
}

func notify(id *idOptions, message string, localize translate) {
	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		return
	}
	defer db.Close()

	var emails []string
	var title string
	if id.User != nil {
		var subscribe bool
		var realname, email string
		if err := db.QueryRow(
			"SELECT subscribe, realname, email FROM user WHERE id = ?",
			id.User).Scan(&subscribe, &realname, &email); err != nil {
			log.Println("Failed to get user subscribe:", err)
		}
		if subscribe {
			title = fmt.Sprintf(localize["Dear"], realname)
			emails = append(emails, email)
		} else {
			return
		}
	} else if len(id.Departments) == 1 {
		var deptName string
		if err := db.QueryRow("SELECT deptname FROM department WHERE id = ?",
			id.Departments[0]).Scan(&deptName); err != nil {
			log.Println("Failed to get department:", err)
			return
		}
		rows, err := db.Query(
			"SELECT email FROM user LEFT JOIN permission p ON id = user_id WHERE subscribe = 1 AND (p.dept_id = ? OR id = 0)",
			id.Departments[0])
		if err != nil {
			log.Println("Failed to get employees:", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var email string
			if err := rows.Scan(&email); err != nil {
				log.Println("Failed to scan email:", err)
				return
			}
			emails = append(emails, email)
		}
		title = fmt.Sprintf(localize["DearAdmin"], deptName)
	} else {
		return
	}
	workers.Slice(emails, func(_ int, email interface{}) {
		if err := Dialer.Send(&mail.Message{
			To:      []string{email.(string)},
			Subject: fmt.Sprintf("%s-%s", localize["NotificationSubject"], time.Now().Format("20060102 15:04")),
			Body:    fmt.Sprintf("%s\n\n    %s", title, message),
		}); err != nil {
			log.Println("Failed to send mail:", err)
		}
	})
}
