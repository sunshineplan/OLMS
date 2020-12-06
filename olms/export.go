package olms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/utils"
)

func (r *record) format(localize map[string]string) (f map[string]interface{}) {
	b, _ := json.Marshal(r)
	json.Unmarshal(b, &f)
	for k, v := range f {
		switch k {
		case "Date":
			f[localize[k]] = strings.Split(f["Date"].(string), "T")[0]
		case "Type":
			switch v.(bool) {
			case false:
				f[localize[k]] = localize["Leave"]
			case true:
				f[localize[k]] = localize["Overtime"]
			}
		case "Status":
			switch int(v.(float64)) {
			case 0:
				f[localize[k]] = localize["Unverified"]
			case 1:
				f[localize[k]] = localize["Verified"]
			case 2:
				f[localize[k]] = localize["Rejected"]
			}
		default:
			f[localize[k]] = v
		}
	}
	return
}

func (s *statistic) format(localize map[string]string) (f map[string]interface{}) {
	b, _ := json.Marshal(s)
	json.Unmarshal(b, &f)
	for k, v := range f {
		f[localize[k]] = v
	}
	return
}

func (d *department) format() map[string]interface{} { return nil }
func (e *employee) format() map[string]interface{}   { return nil }

func sendCSV(c *gin.Context, filename string, fieldnames []string, rows []map[string]interface{}) {
	if len(rows) == 0 {
		c.String(404, "No result.")
		return
	}
	var b bytes.Buffer
	if err := utils.ExportUTF8CSV(fieldnames, rows, &b); err != nil {
		c.String(500, "Failed to save csv: "+err.Error())
		return
	}
	c.Header("content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", url.PathEscape(filename)))
	c.Data(200, "text/csv", b.Bytes())
}
