package olms

import (
	"fmt"
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
		log.Fatalf("Failed to get Self path: %v", err)
	}
	os.MkdirAll(joinPath(dir(Self), "instance"), 0755)
	sqlite = joinPath(dir(Self), "instance", "olms.db")
	sqlitePy = joinPath(dir(Self), "scripts/sqlite.py")
}

func verifyResponse(action, remoteip string, response interface{}) bool {
	if SiteKey != "" && SecretKey != "" {
		if !challenge(action, remoteip, response) {
			return false
		}
	}
	return true
}

func checkSuper(c *gin.Context) bool {
	userID := sessions.Default(c).Get("userID")
	if userID == "0" {
		return true
	}
	return false
}

func checkPermission(c *gin.Context, ids ...interface{}) bool {
	userID := sessions.Default(c).Get("userID")
	if userID == "0" {
		return true
	}
	users, _, err := getEmployees(&idOptions{User: userID}, nil)
	if err != nil {
		return false
	}
	switch len(ids) {
	case 1:
		id := fmt.Sprintf("%v", ids[0])
		for _, i := range strings.Split(users[0].Permission, ",") {
			if id == i {
				return true
			}
		}
	case 2:
		id := fmt.Sprintf("%v", ids[0])
		empls, _, err := getEmployees(&idOptions{User: ids[1]}, nil)
		if err != nil {
			return false
		}
		for _, i := range strings.Split(users[0].Permission, ",") {
			if (ids[0] == nil || id == i) && strconv.Itoa(empls[0].DeptID) == i {
				return true
			}
		}
	}
	return false
}
