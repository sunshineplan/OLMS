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
		db, err := getDB()
		if err != nil {
			log.Println("Failed to connect to database:", err)
			c.AbortWithStatus(503)
			return
		}
		defer db.Close()
		var role bool
		if err := db.QueryRow("SELECT role FROM user WHERE id = ?", userID).Scan(&role); err != nil {
			log.Println("Failed to get role:", err)
			c.AbortWithStatus(500)
			return
		}
		if !role {
			c.AbortWithStatus(403)
		}
	}
}

func superRequired(c *gin.Context) {
	userID := sessions.Default(c).Get("userID")
	if userID == nil {
		c.AbortWithStatus(401)
	} else if userID != "0" {
		c.AbortWithStatus(403)
	}
}

func login(c *gin.Context) {
	localize := localize(c)
	if !verifyResponse("login", c.ClientIP(), c.PostForm("g-recaptcha-response")) {
		c.HTML(200, "login.html", gin.H{"localize": localize, "error": "reCAPTCHA challenge failed", "recaptcha": SiteKey})
		return
	}
	session := sessions.Default(c)
	username := strings.TrimSpace(strings.ToLower(c.PostForm("username")))
	password := c.PostForm("password")

	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		c.String(503, "")
		return
	}
	defer db.Close()
	var id, realname, pw, message string
	if err := db.QueryRow("SELECT id, realname, password FROM user WHERE username = ?",
		username).Scan(&id, &realname, &pw); err != nil {
		if strings.Contains(err.Error(), "no such table") {
			Restore("")
			message = "InitDatabase"
		} else if strings.Contains(err.Error(), "no rows") {
			message = "IncorrectUsername"
		} else {
			log.Print(err)
			message = "CriticalError"
		}
	} else if err = bcrypt.CompareHashAndPassword([]byte(pw), []byte(password)); err != nil {
		if (strings.Contains(err.Error(), "too short") && pw != password) || strings.Contains(err.Error(), "is not") {
			message = "IncorrectPassword"
		} else if pw != password {
			log.Print(err)
			message = "CriticalError"
		}
	}
	if message == "" {
		session.Clear()
		session.Set("userID", id)

		rememberme := c.PostForm("rememberme")
		if rememberme == "on" {
			session.Options(sessions.Options{Path: "/", HttpOnly: true, MaxAge: 856400 * 365})
		} else {
			session.Options(sessions.Options{Path: "/", HttpOnly: true, MaxAge: 0})
		}

		if err := session.Save(); err != nil {
			log.Println("Failed to save session:", err)
			c.String(500, "")
			return
		}
		c.Redirect(302, "/")
		return
	}
	if SiteKey != "" && SecretKey != "" {
		c.HTML(200, "login.html", gin.H{"localize": localize, "error": localize[message], "recaptcha": SiteKey})
		return
	}
	c.HTML(200, "login.html", gin.H{"localize": localize, "error": localize[message]})
}

func setting(c *gin.Context) {
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
		log.Println("Failed to get user subscribe:", err)
	}
	c.HTML(200, "setting.html", gin.H{"localize": localize(c), "subscribe": subscribe, "email": email})
}

func doSetting(c *gin.Context) {
	if !verifyResponse("setting", c.ClientIP(), c.PostForm("g-recaptcha-response")) {
		c.JSON(200, gin.H{"status": 0, "message": "reCAPTCHA challenge failed"})
		return
	}
	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
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
	if err := db.QueryRow("SELECT password FROM user WHERE id = ?", userID).Scan(&oldPassword); err != nil {
		log.Print(err)
		c.String(500, "")
		return
	}

	var message string
	var errorCode int
	err = bcrypt.CompareHashAndPassword([]byte(oldPassword), []byte(password))
	switch {
	case err != nil:
		if (strings.Contains(err.Error(), "too short") && password != oldPassword) || strings.Contains(err.Error(), "is not") {
			message = "IncorrectPassword"
			errorCode = 1
			break
		} else if password != oldPassword {
			log.Print(err)
			c.String(500, "")
			return
		}
		fallthrough
	case password1 == password:
		message = "SamePassword"
		errorCode = 2
	case password1 != password2:
		message = "ConfirmPasswordNoMatch"
		errorCode = 2
	case password1 == "":
		message = "NewPasswordBlank"
	}

	if message == "" {
		newPassword, err := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.MinCost)
		if err != nil {
			log.Print(err)
			c.String(500, "")
			return
		}
		_, err = db.Exec("UPDATE user SET password = ? WHERE id = ?", string(newPassword), userID)
		if err != nil {
			log.Print(err)
			c.String(500, "")
			return
		}
		session.Clear()
		if err := session.Save(); err != nil {
			log.Print(err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		return
	}
	c.JSON(200, gin.H{"status": 0, "message": localize(c)[message], "error": errorCode})
}
