package olms

import (
	"database/sql"
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
	case 0:
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
	} else if userID != 0 {
		c.AbortWithStatus(403)
	}
}

func login(c *gin.Context) {
	var login struct {
		Username, Password string
		Rememberme         bool
		Recaptcha          string
	}
	if err := c.BindJSON(&login); err != nil {
		log.Println("Failed to get option:", err)
		c.String(400, "")
		return
	}
	login.Username = strings.ToLower(login.Username)

	if !verifyResponse("login", c.ClientIP(), login.Recaptcha) {
		c.String(403, "reCAPTCHAChallengeFailed")
		return
	}

	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	var user employee
	statusCode := 200
	var message string
	if err := db.QueryRow(
		"SELECT id, realname, password FROM user WHERE username = ?", login.Username,
	).Scan(&user.ID, &user.Realname, &user.Password); err != nil {
		if strings.Contains(err.Error(), "doesn't exist") {
			Restore("")
			statusCode = 503
			message = "InitDatabase"
		} else if err == sql.ErrNoRows {
			statusCode = 403
			message = "IncorrectUsername"
		} else {
			log.Print(err)
			statusCode = 500
			message = "CriticalError"
		}
	} else {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
			if (err == bcrypt.ErrHashTooShort && user.Password != login.Password) ||
				err == bcrypt.ErrMismatchedHashAndPassword {
				statusCode = 403
				message = "IncorrectPassword"
			} else if user.Password != login.Password {
				log.Print(err)
				statusCode = 500
				message = "CriticalError"
			}
		}
		if message == "" {
			session := sessions.Default(c)
			session.Clear()
			session.Set("userID", user.ID)

			if login.Rememberme {
				session.Options(sessions.Options{Path: "/", HttpOnly: true, MaxAge: 856400 * 365})
			} else {
				session.Options(sessions.Options{Path: "/", HttpOnly: true})
			}

			if err := session.Save(); err != nil {
				log.Print(err)
				statusCode = 500
				message = "Failed to save session."
			}
		}
	}
	c.String(statusCode, message)
}

func setting(c *gin.Context) {
	var setting struct{ Password, Password1, Password2, Recaptcha string }
	if err := c.BindJSON(&setting); err != nil {
		log.Println("Failed to get option:", err)
		c.String(400, "")
		return
	}
	if !verifyResponse("setting", c.ClientIP(), setting.Recaptcha) {
		c.JSON(200, gin.H{"status": 0, "message": "reCAPTCHAChallengeFailed"})
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

	var oldPassword string
	if err = db.QueryRow("SELECT password FROM user WHERE id = ?", userID).Scan(&oldPassword); err != nil {
		log.Print(err)
		c.String(500, "")
		return
	}

	var message string
	var errorCode int
	err = bcrypt.CompareHashAndPassword([]byte(oldPassword), []byte(setting.Password))
	switch {
	case err != nil:
		if (err == bcrypt.ErrHashTooShort && setting.Password != oldPassword) ||
			err == bcrypt.ErrMismatchedHashAndPassword {
			message = "IncorrectPassword"
			errorCode = 1
		} else if setting.Password != oldPassword {
			log.Print(err)
			c.String(500, "")
			return
		}
	case setting.Password1 != setting.Password2:
		message = "ConfirmPasswordNoMatch"
		errorCode = 2
	case setting.Password1 == setting.Password:
		message = "SamePassword"
		errorCode = 2
	case setting.Password1 == "":
		message = "NewPasswordBlank"
	}

	if message == "" {
		newPassword, err := bcrypt.GenerateFromPassword([]byte(setting.Password1), bcrypt.MinCost)
		if err != nil {
			log.Print(err)
			c.String(500, "")
			return
		}
		if _, err := db.Exec("UPDATE user SET password = ? WHERE id = ?", string(newPassword), userID); err != nil {
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
	c.JSON(200, gin.H{"status": 0, "message": message, "error": errorCode})
}
