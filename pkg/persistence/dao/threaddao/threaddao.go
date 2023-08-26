package threaddao

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/orpheus/strings/pkg/infra/sqldb"
)

type ThreadRecord struct {
	Id uuid.UUID
}

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
		return nil, fmt.Errorf("failed to exec thread insert: %s\n", row.Err())
	}

	var threadRecord ThreadRecord
	err := row.Scan(&threadRecord.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to scan thread record: %s\n", err)
	}

	return &threadRecord, nil
}
