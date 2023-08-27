package stringdao

import (
	"fmt"
	"github.com/orpheus/strings/pkg/infra/sqldb"
)

type StringDao struct {
	*sqldb.Store
}

func (s *StringDao) Save(record *StringRecord) (*StringRecord, error) {
	sql := `
	insert into string (id) values ($1) returning id;
	`

	ex := s.Store.GetExecutor()
	row := ex.QueryRow(sql, record.Id)

	if row.Err() != nil {
		return nil, fmt.Errorf("failed to exec string insert: %s", row.Err())
	}

	var r StringRecord
	err := row.Scan(&r.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to scan string record: %s", err)
	}

	return &r, nil
}
