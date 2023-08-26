package threaddao

import (
	"database/sql"
	"errors"
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
	query := `
	insert into thread_versioned (
		id, name, version, thread_id
	) values ($1, $2, $3, $4) 
	returning id, name, version, thread_id, archived, deleted, date_created;
	`

	ex := t.Store.GetExecutor()
	row := ex.QueryRow(query, record.Id, record.Name, record.Version, record.ThreadId)

	var r VersionedThreadRecord
	err := row.Scan(&r.Id, &r.Name, &r.Version, &r.ThreadId, &r.Archived, &r.Deleted, &r.DateCreated)
	if err != nil {
		return nil, fmt.Errorf("Scan err for versioned thread record: %s", err)
	}

	return &r, nil
}

func (t *VersionedThreadDao) FindByThreadId(threadId uuid.UUID) (*VersionedThreadRecord, error) {
	query := `
	select * from thread_versioned where thread_id = $1;
	`

	ex := t.Store.GetExecutor()
	row := ex.QueryRow(query, threadId)

	var r VersionedThreadRecord
	if err := row.Scan(&r.Id, &r.Name, &r.Version, &r.ThreadId, &r.Archived, &r.Deleted, &r.DateCreated); err != nil {
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}
			return nil, fmt.Errorf("QueryRow err for versioned thread record: %s", err)
		}
	}

	return &r, nil
}
