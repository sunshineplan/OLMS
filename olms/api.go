package olms

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func get(c *gin.Context) {
	deptID := c.Query("dept")

	if deptID != "" && !checkPermission(c, deptID) {
		c.String(403, "")
		return
	}
	session := sessions.Default(c)
	users, _, err := getEmpls(session.Get("userID"), nil, nil, nil)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.String(500, "")
		return
	}
	Type := c.Query("type")
	if page, err := strconv.Atoi(c.Query("page")); err != nil {
		fmt.Println(page)
	}
	fmt.Println(users, Type)
	c.JSON(200, nil)
}
