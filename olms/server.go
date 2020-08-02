package olms

import (
	"crypto/rand"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func loadTemplates() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("base.html", joinPath(dir(Self), "templates/base.html"), joinPath(dir(Self), "templates/root.html"))
	r.AddFromFiles("login.html", joinPath(dir(Self), "templates/base.html"), joinPath(dir(Self), "templates/auth/login.html"))
	r.AddFromFiles("setting.html", joinPath(dir(Self), "templates/auth/setting.html"))

	includes, err := filepath.Glob(joinPath(dir(Self), "templates/bookmark/*"))
	if err != nil {
		log.Fatalf("Failed to glob bookmark templates: %v", err)
	}

	for _, include := range includes {
		r.AddFromFiles(filepath.Base(include), include)
	}
	return r
}

// Run server
func Run() {
	f, err := os.OpenFile(LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	gin.DefaultWriter = io.MultiWriter(f)
	log.SetOutput(gin.DefaultWriter)

	secret := make([]byte, 16)
	_, err = rand.Read(secret)
	if err != nil {
		log.Fatalf("Failed to get secret: %v", err)
	}

	router := gin.Default()
	router.Use(sessions.Sessions("session", sessions.NewCookieStore(secret)))
	router.StaticFS("/static", http.Dir(joinPath(dir(Self), "static")))
	router.HTMLRender = loadTemplates()
	router.GET("/", func(c *gin.Context) {
		session := sessions.Default(c)
		name := session.Get("name")
		if name == nil {
			c.Redirect(302, "/auth/login")
			return
		}
		c.HTML(200, "index.html", gin.H{"user": name})
	})
	router.POST("/get", authRequired, get)
	router.POST("/export", authRequired, exportCSV)

	auth := router.Group("/auth")
	auth.GET("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user_id")
		if user != nil {
			c.Redirect(302, "/")
			return
		}
		c.HTML(200, "login.html", gin.H{"error": ""})
	})
	auth.POST("/login", login)
	auth.GET("/logout", authRequired, func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(302, "/auth/login")
	})
	auth.GET("/setting", authRequired, func(c *gin.Context) {
		c.HTML(200, "setting.html", nil)
	})
	auth.POST("/setting", authRequired, setting)

	record := router.Group("/")
	record.Use(authRequired)
	record.GET("/", func(c *gin.Context) {
		c.HTML(200, "showRecords.html", nil)
	})
	record.GET("/add", func(c *gin.Context) {
		c.HTML(200, "addRecord.html", nil)
	})
	record.POST("/add", doAddRecord)
	record.GET("/edit/:id", editRecord)
	record.POST("/edit/:id", doEditRecord)
	record.POST("/delete/:id", doDeleteRecord)

	router.GET("/stats", authRequired, func(c *gin.Context) {
		c.HTML(200, "showStats.html", nil)
	})

	empl := router.Group("/empl")
	empl.GET("/", adminRequired, func(c *gin.Context) {
		c.HTML(200, "showEmpls.html", nil)
	})
	empl.GET("/add", adminRequired, func(c *gin.Context) {
		c.HTML(200, "addEmpl.html", nil)
	})
	empl.POST("/add", adminRequired, doAddEmpl)
	empl.GET("/edit/:id", superRequired, editEmpl)
	empl.POST("/edit/:id", superRequired, doEditEmpl)
	empl.POST("/delete/:id", superRequired, doDeleteEmpl)

	dept := router.Group("/dept")
	dept.Use(superRequired)
	dept.GET("/", func(c *gin.Context) {
		c.HTML(200, "showDepts.html", nil)
	})
	dept.GET("/add", func(c *gin.Context) {
		c.HTML(200, "addDept.html", nil)
	})
	dept.POST("/add", doAddDept)
	dept.GET("/edit/:id", editDept)
	dept.POST("/edit/:id", doEditDept)
	dept.POST("/delete/:id", doDeleteDept)

	if UNIX != "" && OS == "linux" {
		if _, err := os.Stat(UNIX); err == nil {
			err = os.Remove(UNIX)
			if err != nil {
				log.Fatalf("Failed to remove socket file: %v", err)
			}
		}

		listener, err := net.Listen("UNIX", UNIX)
		if err != nil {
			log.Fatalf("Failed to listen socket file: %v", err)
		}

		idleConnsClosed := make(chan struct{})
		go func() {
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit

			if err := listener.Close(); err != nil {
				log.Printf("Failed to close listener: %v", err)
			}
			if _, err := os.Stat(UNIX); err == nil {
				err = os.Remove(UNIX)
				if err != nil {
					log.Printf("Failed to remove socket file: %v", err)
				}
			}
			close(idleConnsClosed)
		}()

		if err = os.Chmod(UNIX, 0666); err != nil {
			log.Fatalf("Failed to chmod socket file: %v", err)
		}

		http.Serve(listener, router)
		<-idleConnsClosed
	} else {
		router.Run(Host + ":" + Port)
	}
}
