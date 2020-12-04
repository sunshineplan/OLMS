package olms

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Run server
func Run() {
	if LogPath != "" {
		f, err := os.OpenFile(LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		gin.DefaultWriter = f
		gin.DefaultErrorWriter = f
		log.SetOutput(f)
	}

	secret := make([]byte, 16)
	if _, err := rand.Read(secret); err != nil {
		log.Fatalf("Failed to get secret: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	Server.Handler = router
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.Use(sessions.Sessions("session", sessions.NewCookieStore(secret)))
	router.StaticFS("/js", http.Dir(joinPath(dir(Self), "dist/js")))
	router.StaticFS("/css", http.Dir(joinPath(dir(Self), "dist/css")))
	router.StaticFile("favicon.ico", joinPath(dir(Self), "dist/favicon.ico"))
	router.LoadHTMLFiles(joinPath(dir(Self), "dist/index.html"))
	router.GET("/", func(c *gin.Context) {
		var user employee
		switch userID := sessions.Default(c).Get("userID"); userID {
		case nil:
			c.Redirect(302, "/auth/login")
		case "0":
			user = employee{ID: 0, Realname: "root", Role: true}
		default:
			users, _, err := getEmployees(&idOptions{User: userID}, nil)
			if err != nil {
				log.Printf("Failed to get users: %v", err)
				c.String(500, "")
				return
			}
			user = users[0]
		}
		if SiteKey != "" && SecretKey != "" {
			c.HTML(200, "index.html", gin.H{"localize": localize(c), "user": user, "recaptcha": SiteKey})
			return
		}
		c.HTML(200, "index.html", gin.H{"localize": localize(c), "user": user})
	})

	auth := router.Group("/")
	auth.GET("/login", func(c *gin.Context) {
		user := sessions.Default(c).Get("userID")
		if user != nil {
			c.Redirect(302, "/")
			return
		}
		if SiteKey != "" && SecretKey != "" {
			c.HTML(200, "login.html", gin.H{"localize": localize(c), "error": "", "recaptcha": SiteKey})
			return
		}
		c.HTML(200, "login.html", gin.H{"localize": localize(c), "error": ""})
	})
	auth.POST("/login", login)
	auth.GET("/logout", authRequired, func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(302, "/login")
	})
	auth.GET("/setting", authRequired, setting)
	auth.POST("/setting", authRequired, doSetting)

	api := router.Group("/")
	api.Use(authRequired)
	api.POST("/get", get)
	api.POST("/export", exportCSV)
	api.POST("/subscribe", subscribe)

	record := router.Group("/record")
	record.Use(authRequired)
	record.POST("/add", addRecord)
	record.POST("/edit", editRecord)
	record.POST("/verify", adminRequired, verifyRecord)
	record.POST("/delete", deleteRecord)

	empl := router.Group("/employee")
	empl.POST("/add", adminRequired, addEmployee)
	empl.POST("/edit", superRequired, editEmployee)
	empl.POST("/delete", superRequired, deleteEmployee)

	dept := router.Group("/department")
	dept.Use(superRequired)
	dept.POST("/add", addDepartment)
	dept.POST("/edit", editDepartment)
	dept.POST("/delete", deleteDepartment)

	router.NoRoute(func(c *gin.Context) {
		c.Redirect(302, "/")
	})

	if err := Server.Run(); err != nil {
		log.Fatal(err)
	}
}
