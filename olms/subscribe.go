package olms

import (
	"log"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func subscribe(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	userID := sessions.Default(c).Get("userID")
	switch c.Query("subscribe") {
	case "0":
		if _, err := db.Exec("UPDATE user SET subscribe = 0 WHERE id = ?", userID); err != nil {
			log.Printf("Failed to update user: %v", err)
			c.String(500, "")
			return
		}
	case "1":
		if _, err := db.Exec(
			"UPDATE user SET subscribe = 1 AND email = ? WHERE id = ?", c.Query("email"), userID); err != nil {
			log.Printf("Failed to update user: %v", err)
			c.String(500, "")
			return
		}
	default:
		c.String(400, "")
		return
	}
	c.JSON(200, gin.H{"status": 1})
}

func notify(id *idOptions, template string) {
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	if id.UserID != nil {
		var ok bool
		var realname, email string
		if err := db.QueryRow(
			"SELECT subscribe, realname, email FROM user WHERE id = ?", id).Scan(&ok, &realname, &email); err != nil {
			log.Printf("Failed to get user subscribe: %v", err)
		}
		if ok {
			log.Printf("mail to %s %s: %s", email, realname, template)
		}
	}
	if len(id.DeptIDs) == 1 {
		var deptName string
		if err := db.QueryRow("SELECT dept_name FROM department WHERE id = ?", id).Scan(&deptName); err != nil {
			log.Printf("Failed to get department: %v", err)
			return
		}
		rows, err := db.Query(
			"SELECT email FROM user JOIN permission p ON id = user_id WHERE subscribe = 1 AND p.dept_id = ?", id.DeptIDs[0])
		if err != nil {
			log.Printf("Failed to get employees: %v", err)
			return
		}
		defer rows.Close()
		var emails []string
		for rows.Next() {
			var email string
			if err = rows.Scan(&email); err != nil {
				log.Printf("Failed to scan email: %v", err)
				return
			}
			emails = append(emails, email)
		}
		log.Printf("mail to %v %s: %s", emails, deptName, template)
	}
}
