package olms

import (
	"log"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func authRequired(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		c.AbortWithStatus(401)
	}
}

func adminRequired(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	switch userID {
	case nil:
		c.AbortWithStatus(401)
	case "0":
	default:
		user, _, err := getEmpls(userID, nil, nil, nil)
		if err != nil {
			c.AbortWithStatus(500)
		} else if !user[0].Role {
			c.AbortWithStatus(403)
		}
	}
}

func superRequired(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		c.AbortWithStatus(401)
	} else if userID != "0" {
		c.AbortWithStatus(403)
	}
}

func login(c *gin.Context) {
	session := sessions.Default(c)
	username := strings.TrimSpace(strings.ToLower(c.PostForm("username")))
	password := c.PostForm("password")

	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.HTML(200, "login.html", gin.H{"error": "Failed to connect to database."})
		return
	}
	defer db.Close()
	var id, realname, pw, message string
	if err := db.QueryRow("SELECT id, realname, password FROM user WHERE username = ?",
		username).Scan(&id, &realname, &pw); err != nil {
		if strings.Contains(err.Error(), "no such table") {
			Restore("")
			c.HTML(200, "login.html", gin.H{"error": "Detected first time running. Initialized the database."})
			return
		}
		if strings.Contains(err.Error(), "no rows") {
			message = "Incorrect username"
		} else {
			log.Println(err)
			c.HTML(200, "login.html", gin.H{"error": "Critical Error! Please contact your system administrator."})
			return
		}
	} else {
		if err = bcrypt.CompareHashAndPassword([]byte(pw), []byte(password)); err != nil {
			if (strings.Contains(err.Error(), "too short") && pw != password) || strings.Contains(err.Error(), "is not") {
				message = "Incorrect password"
			} else if pw != password {
				log.Println(err)
				c.HTML(200, "login.html", gin.H{"error": "Critical Error! Please contact your system administrator."})
				return
			}
		}
		if message == "" {
			session.Clear()
			session.Set("userID", id)
			session.Set("name", realname)

			rememberme := c.PostForm("rememberme")
			if rememberme == "on" {
				session.Options(sessions.Options{Path: "/", HttpOnly: true, MaxAge: 856400 * 365})
			} else {
				session.Options(sessions.Options{Path: "/", HttpOnly: true, MaxAge: 0})
			}

			if err := session.Save(); err != nil {
				log.Println(err)
				c.HTML(200, "login.html", gin.H{"error": "Failed to save session."})
				return
			}
			c.Redirect(302, "/")
			return
		}
	}
	c.HTML(200, "login.html", gin.H{"error": message})
}

func setting(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.String(503, "")
		return
	}
	defer db.Close()
	session := sessions.Default(c)
	userID := session.Get("userID")

	password := c.PostForm("password")
	password1 := c.PostForm("password1")
	password2 := c.PostForm("password2")

	var oldPassword string
	err = db.QueryRow("SELECT password FROM user WHERE id = ?", userID).Scan(&oldPassword)
	if err != nil {
		log.Println(err)
		c.String(500, "")
		return
	}

	var message string
	var errorCode int
	err = bcrypt.CompareHashAndPassword([]byte(oldPassword), []byte(password))
	switch {
	case err != nil:
		if (strings.Contains(err.Error(), "too short") && password != oldPassword) || strings.Contains(err.Error(), "is not") {
			message = "Incorrect password."
			errorCode = 1
		} else if password != oldPassword {
			log.Println(err)
			c.String(500, "")
			return
		}
	case password1 != password2:
		message = "Confirm password doesn't match new password."
		errorCode = 2
	case password1 == password:
		message = "New password cannot be the same as your current password."
		errorCode = 2
	case password1 == "":
		message = "New password cannot be blank."
	}

	if message == "" {
		newPassword, err := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.MinCost)
		if err != nil {
			log.Println(err)
			c.String(500, "")
			return
		}
		_, err = db.Exec("UPDATE user SET password = ? WHERE id = ?", string(newPassword), userID)
		if err != nil {
			log.Println(err)
			c.String(500, "")
			return
		}
		session.Clear()
		if err := session.Save(); err != nil {
			log.Println(err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		return
	}
	c.JSON(200, gin.H{"status": 0, "message": message, "error": errorCode})
}
