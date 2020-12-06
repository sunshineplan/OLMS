package olms

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/utils/httpsvr"
)

// Self execute file location
var Self string

// Server is an HTTP server
var Server httpsvr.Server

// LogPath log file location
var LogPath string

var perPage = 10
var sqlite, sqlitePy string

func init() {
	var err error
	Self, err = os.Executable()
	if err != nil {
		log.Fatalln("Failed to get Self path:", err)
	}
	os.MkdirAll(joinPath(dir(Self), "instance"), 0755)
	sqlite = joinPath(dir(Self), "instance", "olms.db")
	sqlitePy = joinPath(dir(Self), "scripts/sqlite.py")
}

func verifyResponse(action, remoteip, response string) bool {
	if SiteKey != "" && SecretKey != "" {
		if !challenge(action, remoteip, response) {
			return false
		}
	}
	return true
}

func checkPermission(db *sql.DB, c *gin.Context, option *idOptions) bool {
	userID := sessions.Default(c).Get("userID")
	if userID == "0" {
		return true
	}

	user, err := getUser(db, userID)
	if err != nil {
		return false
	}
	switch {
	case option.Departments != nil && option.User == nil:
		for _, i := range strings.Split(user.Permission, ",") {
			if option.Departments[0] == i {
				return true
			}
		}
	default:
		employee, err := getUser(db, option.User)
		if err != nil {
			return false
		}
		for _, i := range strings.Split(user.Permission, ",") {
			if (option.Departments == nil || option.Departments[0] == i) &&
				strconv.Itoa(employee.DeptID) == i {
				return true
			}
		}
	}
	return false
}

func checkRecord(db *sql.DB, c *gin.Context, id int, self bool) (record record, ok bool) {
	if err := db.QueryRow(
		"SELECT record.dept_id, user_id, realname, status FROM record JOIN user ON user_id = user.id WHERE record.id = ?",
		id).Scan(&record.DeptID, &record.UserID, &record.Realname, &record.Status); err != nil {
		log.Println("Failed to get record:", err)
		return
	}
	if self {
		userID := sessions.Default(c).Get("userID")
		if userID == "0" {
			ok = true
			return
		}
		ok = userID == strconv.Itoa(record.UserID)
		return
	}
	ok = checkPermission(db, c, &idOptions{Departments: []string{strconv.Itoa(record.DeptID)}})
	return
}
