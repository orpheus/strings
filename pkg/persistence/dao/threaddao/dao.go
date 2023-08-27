package threaddao

import (
	"fmt"
	"github.com/orpheus/strings/pkg/infra/sqldb"
)

type ThreadDao struct {
	*sqldb.Store
}

func (t *ThreadDao) Save(record *ThreadRecord) (*ThreadRecord, error) {
	sql := `
	insert into thread (id) values ($1) returning id;
	`

	ex := t.Store.GetExecutor()
	row := ex.QueryRow(sql, record.Id)

	if row.Err() != nil {
		return nil, fmt.Errorf("failed to exec thread insert: %s", row.Err())
	}

	var threadRecord ThreadRecord
	err := row.Scan(&threadRecord.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to scan thread record: %s", err)
	}

	return &threadRecord, nil
}
