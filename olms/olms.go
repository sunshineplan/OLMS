package olms

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
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
	os.MkdirAll(joinPath(dir(Self), "instance"), 0755)
	sqlite = joinPath(dir(Self), "instance", "olms.db")
	sqlitePy = joinPath(dir(Self), "scripts/sqlite.py")
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
	users, _, err := getEmpls(userID, nil, nil, nil)
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
		empls, _, err := getEmpls(ids[1], nil, nil, nil)
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
