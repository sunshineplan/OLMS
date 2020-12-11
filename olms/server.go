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
			log.Fatalln("Failed to open log file:", err)
		}
		gin.DefaultWriter = f
		gin.DefaultErrorWriter = f
		log.SetOutput(f)
	}

	secret := make([]byte, 16)
	if _, err := rand.Read(secret); err != nil {
		log.Fatalln("Failed to get secret:", err)
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
		c.HTML(200, "index.html", nil)
	})

	auth := router.Group("/")
	auth.GET("/info", info)
	auth.POST("/login", login)
	auth.GET("/logout", authRequired, func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(302, "/")
	})
	auth.POST("/setting", authRequired, setting)

	api := router.Group("/")
	api.Use(authRequired)
	api.GET("/year", year)
	api.POST("/year", year)
	api.POST("/records", records)
	api.POST("/statistics", statistics)
	api.POST("/records/export", exportRecords)
	api.POST("/statistics/export", exportStatistics)
	api.GET("/subscribe", getSubscribe)
	api.POST("/subscribe", subscribe)

	record := router.Group("/record")
	record.Use(authRequired)
	record.POST("/add", addRecord)
	record.POST("/edit", editRecord)
	record.POST("/verify/:id", adminRequired, verifyRecord)
	record.POST("/delete/:id", deleteRecord)

	employee := router.Group("/employee")
	employee.POST("/add", adminRequired, addEmployee)
	employee.POST("/edit", superRequired, editEmployee)
	employee.POST("/delete/:id", superRequired, deleteEmployee)

	department := router.Group("/department")
	department.Use(superRequired)
	department.POST("/add", addDepartment)
	department.POST("/edit", editDepartment)
	department.POST("/delete/:id", deleteDepartment)

	router.NoRoute(func(c *gin.Context) {
		c.Redirect(302, "/")
	})

	if err := Server.Run(); err != nil {
		log.Fatal(err)
	}
}
