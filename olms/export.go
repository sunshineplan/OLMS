package olms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/utils/export"
)

func (r record) format() (f map[string]interface{}) {
	b, _ := json.Marshal(r)
	json.Unmarshal(b, &f)
	for k, v := range f {
		switch k {
		case "Type":
			switch v {
			case 0:
				v = "Leave"
			case 1:
				v = "Overtime"
			}
		case "Status":
			switch v {
			case 0:
				v = "Unverified"
			case 1:
				v = "Verified"
			case 2:
				v = "Rejected"
			}
		}
	}
	return
}

func (s stat) format() (f map[string]interface{}) {
	b, _ := json.Marshal(s)
	json.Unmarshal(b, &f)
	return
}

func (d dept) format() map[string]interface{} { return nil }
func (e empl) format() map[string]interface{} { return nil }

func sendCSV(c *gin.Context, filename string, fieldnames []string, r []map[string]interface{}) {
	var rows []interface{}
	for _, i := range r {
		rows = append(rows, i)
	}
	var b bytes.Buffer
	if err := export.CSV(fieldnames, rows, &b); err != nil {
		c.String(500, "Failed to save csv: "+err.Error())
		return
	}
	body, err := ioutil.ReadAll(&b)
	if err != nil {
		c.String(500, "Failed to read csv bytes: "+err.Error())
		return
	}
	c.Header("content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(200, "text/csv", body)
}
