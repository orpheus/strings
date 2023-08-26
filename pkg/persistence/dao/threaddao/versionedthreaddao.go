package threaddao

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/orpheus/strings/pkg/infra/sqldb"
	"time"
)

type VersionedThreadRecord struct {
	Id          uuid.UUID
	Name        string
	Version     int
	ThreadId    uuid.UUID
	Archived    bool
	Deleted     bool
	DateCreated time.Time
}

type VersionedThreadDao struct {
	*sqldb.Store
}

func (t *VersionedThreadDao) Save(record *VersionedThreadRecord) (*VersionedThreadRecord, error) {
	sql := `
	insert into thread_versioned (
		id, name, version, thread_id
	) values ($1, $2, $3, $4) 
	returning id, name, version, thread_id, archived, deleted, date_created;
	`

	ex := t.Store.GetExecutor()
	row := ex.QueryRow(sql, record.Id, record.Name, record.Version, record.ThreadId)

	if row.Err() != nil {
		return nil, fmt.Errorf("failed to exec thread insert: %s\n", row.Err())
	}

	var r VersionedThreadRecord
	err := row.Scan(&r.Id, &r.Name, &r.Version, &r.ThreadId, &r.Archived, &r.Deleted, &r.DateCreated)
	if err != nil {
		return nil, fmt.Errorf("failed to scan versioned thread record: %s\n", err)
	}

	return &r, nil
}
