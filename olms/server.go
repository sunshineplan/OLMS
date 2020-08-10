package olms

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func loadTemplates() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index.html", joinPath(dir(Self), "templates/base.html"), joinPath(dir(Self), "templates/index.html"))
	r.AddFromFiles("login.html", joinPath(dir(Self), "templates/base.html"), joinPath(dir(Self), "templates/auth/login.html"))
	r.AddFromFiles("setting.html", joinPath(dir(Self), "templates/auth/setting.html"))

	dept, err := filepath.Glob(joinPath(dir(Self), "templates/dept/*"))
	if err != nil {
		log.Fatalf("Failed to glob dept templates: %v", err)
	}
	empl, err := filepath.Glob(joinPath(dir(Self), "templates/empl/*"))
	if err != nil {
		log.Fatalf("Failed to glob empl templates: %v", err)
	}
	record, err := filepath.Glob(joinPath(dir(Self), "templates/record/*"))
	if err != nil {
		log.Fatalf("Failed to glob record templates: %v", err)
	}
	stat, err := filepath.Glob(joinPath(dir(Self), "templates/stat/*"))
	if err != nil {
		log.Fatalf("Failed to glob stat templates: %v", err)
	}
	var includes []string
	for _, i := range [][]string{dept, empl, record, stat} {
		includes = append(includes, i...)
	}
	for _, include := range includes {
		r.AddFromFiles(filepath.Base(include), include)
	}
	return r
}

// Run server
func Run() {
	if LogPath != "" {
		f, err := os.OpenFile(LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
		log.SetOutput(gin.DefaultWriter)
	}

	secret := make([]byte, 16)
	if _, err := rand.Read(secret); err != nil {
		log.Fatalf("Failed to get secret: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
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
	router.StaticFS("/static", http.Dir(joinPath(dir(Self), "static")))
	router.HTMLRender = loadTemplates()
	router.GET("/", func(c *gin.Context) {
		switch userID := sessions.Default(c).Get("userID"); userID {
		case nil:
			c.Redirect(302, "/auth/login")
		case "0":
			c.HTML(200, "index.html", gin.H{"user": empl{ID: 0, Realname: "root", Role: true}})
		default:
			users, _, err := getEmpls(userID, nil, "", "")
			if err != nil {
				log.Printf("Failed to get users: %v", err)
				c.String(500, "")
				return
			}
			c.HTML(200, "index.html", gin.H{"user": users[0]})
		}
	})

	auth := router.Group("/auth")
	auth.GET("/login", func(c *gin.Context) {
		user := sessions.Default(c).Get("userID")
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

	api := router.Group("/")
	api.Use(authRequired)
	api.POST("/get", get)
	api.POST("/export", exportCSV)

	record := router.Group("/record")
	record.Use(authRequired)
	record.GET("", func(c *gin.Context) {
		c.HTML(200, "showRecords.html", gin.H{"mode": ""})
	})
	record.GET("/admin", adminRequired, func(c *gin.Context) {
		c.HTML(200, "showRecords.html", gin.H{"mode": "admin"})
	})
	record.GET("/super", superRequired, func(c *gin.Context) {
		c.HTML(200, "showRecords.html", gin.H{"mode": "super"})
	})
	record.GET("/add", func(c *gin.Context) {
		c.HTML(200, "record.html", gin.H{"mode": "", "id": "0", "user": "1"})
	})
	record.GET("/admin/add", adminRequired, func(c *gin.Context) {
		c.HTML(200, "record.html", gin.H{"mode": "admin", "id": "0", "user": sessions.Default(c).Get("userID")})
	})
	record.POST("/add", doAddRecord)
	record.GET("/edit/:id", func(c *gin.Context) {
		c.HTML(200, "record.html", gin.H{"mode": "", "id": c.Param("id"), "user": "1"})
	})
	record.GET("/super/edit/:id", superRequired, func(c *gin.Context) {
		c.HTML(200, "record.html", gin.H{"mode": "super", "id": c.Param("id"), "user": sessions.Default(c).Get("userID")})
	})
	record.POST("/edit/:id", doEditRecord)
	record.GET("/verify/:id", adminRequired, func(c *gin.Context) {
		c.HTML(200, "verify.html", gin.H{"id": c.Param("id")})
	})
	record.POST("/verify/:id", adminRequired, doVerifyRecord)
	record.POST("/delete/:id", doDeleteRecord)

	stat := router.Group("/stat")
	stat.GET("", authRequired, func(c *gin.Context) {
		c.HTML(200, "showStats.html", gin.H{"mode": ""})
	})
	stat.GET("/admin", adminRequired, func(c *gin.Context) {
		c.HTML(200, "showStats.html", gin.H{"mode": "admin"})
	})

	empl := router.Group("/empl")
	empl.GET("", adminRequired, func(c *gin.Context) {
		c.HTML(200, "showEmpls.html", gin.H{"user": sessions.Default(c).Get("userID")})
	})
	empl.GET("/add", adminRequired, func(c *gin.Context) {
		c.HTML(200, "empl.html", gin.H{"id": "0", "user": sessions.Default(c).Get("userID")})
	})
	empl.POST("/add", adminRequired, doAddEmpl)
	empl.GET("/edit/:id", superRequired, func(c *gin.Context) {
		c.HTML(200, "empl.html", gin.H{"id": c.Param("id"), "user": sessions.Default(c).Get("userID")})
	})
	empl.POST("/edit/:id", superRequired, doEditEmpl)
	empl.POST("/delete/:id", superRequired, doDeleteEmpl)

	dept := router.Group("/dept")
	dept.Use(superRequired)
	dept.GET("", func(c *gin.Context) {
		c.HTML(200, "showDepts.html", nil)
	})
	dept.GET("/add", func(c *gin.Context) {
		c.HTML(200, "dept.html", gin.H{"id": "0"})
	})
	dept.POST("/add", doAddDept)
	dept.GET("/edit/:id", func(c *gin.Context) {
		c.HTML(200, "dept.html", gin.H{"id": c.Param("id")})
	})
	dept.POST("/edit/:id", doEditDept)
	dept.POST("/delete/:id", doDeleteDept)

	if UNIX != "" && OS == "linux" {
		if _, err := os.Stat(UNIX); err == nil {
			if err := os.Remove(UNIX); err != nil {
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
				if err := os.Remove(UNIX); err != nil {
					log.Printf("Failed to remove socket file: %v", err)
				}
			}
			close(idleConnsClosed)
		}()

		if err := os.Chmod(UNIX, 0666); err != nil {
			log.Fatalf("Failed to chmod socket file: %v", err)
		}

		http.Serve(listener, router)
		<-idleConnsClosed
	} else {
		router.Run(Host + ":" + Port)
	}
}
