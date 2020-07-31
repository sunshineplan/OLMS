package olms

import (
	"fmt"
	"log"
	"strings"
)

type stat struct {
	Period   string
	DeptName string
	Name     string
	Overtime int
	Leave    int
	Summary  int
}

func getStats(id interface{}, deptIDs []interface{}, period, year, month string, page interface{}) (stats []stat, total int, err error) {
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	stmt := "SELECT %s FROM statistics WHERE "
	var fields, group, order, limit string
	var args []interface{}
	if period == "month" {
		fields = "period, dept_name, realname, overtime, leave, summary"
		if month == "" {
			if year != "" {
				stmt += "substr(period,1,4) = ? AND "
				args = append(args, year)
			}
		} else {
			stmt += "period = ? AND "
			args = append(args, year+"-"+month)
		}
	} else {
		fields = "substr(period,1,4) year, dept_name, realname, sum(overtime), sum(leave), sum(summary)"
		group = " GROUP BY year, dept_id, user_id"
		order = " ORDER BY year DESC"
	}

	if id != nil {
		stmt += " user_id = ?"
		args = append(args, id)
	} else {
		marks := make([]string, len(deptIDs))
		for i := range marks {
			marks[i] = "?"
		}
		stmt += " dept_id IN (" + strings.Join(marks, ", ") + ")"
		args = append(args, deptIDs...)
	}

	bc := make(chan bool, 1)
	if p, ok := page.(int); ok {
		limit = fmt.Sprintf(" LIMIT ?, ?")
		args = append(args, (p-1)*perPage, perPage)
		go func() {
			if err := db.QueryRow(fmt.Sprintf(stmt+group, "count(realname)")).Scan(&total); err != nil {
				log.Printf("Failed to get total records: %v", err)
				bc <- false
			}
			bc <- true
		}()
	} else {
		bc <- true
	}
	rows, err := db.Query(fmt.Sprintf(stmt+group+order+limit, fields), args...)
	if err != nil {
		log.Printf("Failed to get statistics: %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var stat stat
		if err = rows.Scan(&stat.Period, &stat.DeptName, &stat.Name, &stat.Overtime, &stat.Leave, &stat.Summary); err != nil {
			log.Printf("Failed to scan statistics: %v", err)
			return
		}
		stats = append(stats, stat)
	}
	if v := <-bc; !v {
		err = fmt.Errorf("Failed to get total records")
	}
	return
}
