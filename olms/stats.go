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

func getStats(id interface{}, deptIDs []string, period, year, month, page, sort, order interface{}) (stats []stat, total int, err error) {
	db, err := getDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	stmt := "SELECT %s FROM statistics WHERE "
	var fields, group, orderBy, limit string
	var args []interface{}
	if period == "month" {
		fields = "period, dept_name, realname, overtime, leave, summary"
		if month == nil {
			if year != nil {
				stmt += "substr(period,1,4) = ? AND "
				args = append(args, year)
			}
		} else {
			stmt += "period = ? AND "
			args = append(args, fmt.Sprintf("%v-%v", year, month))
		}
	} else {
		fields = "substr(period,1,4) period, dept_name, realname, sum(overtime), sum(leave), sum(summary)"
		group = " GROUP BY period, dept_id, user_id"
		orderBy = " ORDER BY period DESC"
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
		for _, i := range deptIDs {
			args = append(args, i)
		}
	}

	bc := make(chan bool, 1)
	if p, ok := page.(float64); ok {
		go func() {
			if err := db.QueryRow(fmt.Sprintf("SELECT count(*) FROM (%s)",
				fmt.Sprintf(stmt+group, "substr(period,1,4) period")), args...).Scan(&total); err != nil {
				log.Printf("Failed to get total records: %v", err)
				bc <- false
			}
			bc <- true
		}()
		limit = fmt.Sprintf(" LIMIT ?, ?")
		args = append(args, int(p-1)*perPage, perPage)
	} else {
		bc <- true
	}
	if sort != nil {
		orderBy = fmt.Sprintf(" ORDER BY %v %v", sort, order)
	}
	rows, err := db.Query(fmt.Sprintf(stmt+group+orderBy+limit, fields), args...)
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
