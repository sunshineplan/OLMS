package olms

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type statistic struct {
	Period   string `json:"period"`
	DeptName string `json:"deptname"`
	Realname string `json:"realname"`
	Overtime int    `json:"overtime"`
	Leave    int    `json:"leave"`
	Summary  int    `json:"summary"`
}

func getStatistics(db *sql.DB, id *idOptions, options *searchOptions) (statistics []statistic, total int, err error) {
	stmt := "SELECT %s FROM statistics WHERE"

	var args []interface{}
	var fields, group, orderBy, limit string
	bc := make(chan bool, 1)
	if id.User != nil {
		stmt += " user_id = ?"
		args = append(args, id.User)
	} else {
		marks := make([]string, len(id.Departments))
		for i := range marks {
			marks[i] = "?"
		}
		stmt += " dept_id IN (" + strings.Join(marks, ", ") + ")"
		for _, i := range id.Departments {
			args = append(args, i)
		}
	}

	if options.Period == "month" {
		fields = "period, deptname, realname, overtime, leave, summary"
		if options.Month == "" {
			if options.Year != "" {
				stmt += " AND substr(period,1,4) = ?"
				args = append(args, options.Year)
			}
		} else {
			stmt += " AND period = ?"
			args = append(args, fmt.Sprintf("%v-%v", options.Year, options.Month))
		}
	} else {
		fields = "substr(period,1,4) period, deptname, realname, sum(overtime), sum(leave), sum(summary)"
		group = " GROUP BY period, dept_id, user_id"
		orderBy = " ORDER BY period DESC"
	}
	if options.Page != 0 {
		go func() {
			if err := db.QueryRow(fmt.Sprintf("SELECT count(*) FROM (%s)",
				fmt.Sprintf(stmt+group, "substr(period,1,4) period")), args...).Scan(&total); err != nil {
				log.Println("Failed to get total records:", err)
				bc <- false
			}
			bc <- true
		}()
		limit = fmt.Sprintf(" LIMIT ?, ?")
		args = append(args, int(options.Page-1)*perPage, perPage)
	} else {
		bc <- true
	}
	if options.Sort != "" {
		orderBy = fmt.Sprintf(" ORDER BY %v %v", options.Sort, options.Order)
	}
	rows, err := db.Query(fmt.Sprintf(stmt+group+orderBy+limit, fields), args...)
	if err != nil {
		log.Println("Failed to get statistics:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var statistic statistic
		if err = rows.Scan(
			&statistic.Period, &statistic.DeptName, &statistic.Realname, &statistic.Overtime, &statistic.Leave, &statistic.Summary,
		); err != nil {
			log.Println("Failed to scan statistics:", err)
			return
		}
		statistics = append(statistics, statistic)
	}
	if v := <-bc; !v {
		err = fmt.Errorf("Failed to get total records")
	}
	return
}
