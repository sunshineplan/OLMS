package olms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/utils/export"
)

func (r record) format() (f map[string]interface{}) {
	b, _ := json.Marshal(r)
	json.Unmarshal(b, &f)
	f["Date"] = strings.Split(f["Date"].(string), "T")[0]
	for k, v := range f {
		switch k {
		case "Type":
			switch v.(bool) {
			case false:
				f[k] = "Leave"
			case true:
				f[k] = "Overtime"
			}
		case "Status":
			switch int(v.(float64)) {
			case 0:
				f[k] = "Unverified"
			case 1:
				f[k] = "Verified"
			case 2:
				f[k] = "Rejected"
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
	if len(r) == 0 {
		c.String(404, "No result.")
		return
	}
	var rows []interface{}
	for _, i := range r {
		rows = append(rows, i)
	}
	var b bytes.Buffer
	b.Write([]byte{0xEF, 0xBB, 0xBF})
	if err := export.CSV(fieldnames, rows, &b); err != nil {
		c.String(500, "Failed to save csv: "+err.Error())
		return
	}
	body, err := ioutil.ReadAll(&b)
	if err != nil {
		c.String(500, "Failed to read csv bytes: "+err.Error())
		return
	}
	c.Header("content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", url.PathEscape(filename)))
	c.Data(200, "text/csv", body)
}
