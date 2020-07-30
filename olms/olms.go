package olms

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// OS is the running program's operating system
const OS = runtime.GOOS

// Self execute file location
var Self string

// UNIX file
var UNIX string

// Host address
var Host string

// Port number
var Port string

// LogPath log file location
var LogPath string

var perPage = 10
var sqlite, sqlitePy string

func init() {
	var err error
	Self, err = os.Executable()
	if err != nil {
		log.Fatalf("Failed to get Self path: %v", err)
	}
	os.MkdirAll(filepath.Join(filepath.Dir(Self), "instance"), 0755)
	sqlite = filepath.Join(filepath.Dir(Self), "instance", "olms.db")
	sqlitePy = filepath.Join(filepath.Dir(Self), "scripts/sqlite.py")
}

func authRequired(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		c.AbortWithStatus(401)
	}
}

func adminRequired(c *gin.Context) {
	session := sessions.Default(c)
	user, err := getEmpl(session.Get("userID"))
	if err != nil {
		c.AbortWithStatus(401)
	} else if !user.Admin {
		c.AbortWithStatus(403)
	}
}

func superRequired(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		c.AbortWithStatus(401)
	} else if userID != 0 {
		c.AbortWithStatus(403)
	}
}

func getDept(id interface{}) (dept, error) {
	var dept dept
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return dept, err
	}
	defer db.Close()
	if err := db.QueryRow("SELECT * FROM department WHERE id = ?", id).Scan(&dept.ID, &dept.Name); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Printf("No results: %v", err)
		} else {
			log.Printf("Failed to query dept: %v", err)
		}
		return dept, err
	}
	return dept, nil
}

func getEmpl(id interface{}) (empl, error) {
	var empl empl
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return empl, err
	}
	defer db.Close()
	if err := db.QueryRow("SELECT id, realname, dept_id, dept_name, admin, permission FROM employee WHERE id = ?",
		id).Scan(&empl.ID, &empl.Realname, &empl.DeptID, &empl.DeptName, &empl.Admin, &empl.Permission); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Printf("No results: %v", err)
		} else {
			log.Printf("Failed to query employee: %v", err)
		}
		return empl, err
	}
	return empl, nil
}

func checkPermission(id string, c *gin.Context) bool {
	session := sessions.Default(c)
	user, err := getEmpl(session.Get("userID"))
	if err != nil {
		return false
	}
	for _, i := range strings.Split(user.Permission, ",") {
		if id == i {
			return true
		}
	}
	return false
}
