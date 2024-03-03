package zinc

import (
	"database/sql"
)

type Result struct {
	LastInsertId int64
	RowsAffected int64
}

func newResult(sqlRes sql.Result) (*Result, error) {
	var execRes Result
	if n, err := sqlRes.LastInsertId(); err != nil {
		return nil, err
	} else {
		execRes.LastInsertId = n
	}
	if n, err := sqlRes.RowsAffected(); err != nil {
		return nil, err
	} else {
		execRes.RowsAffected = n
	}
	return &execRes, nil
}
