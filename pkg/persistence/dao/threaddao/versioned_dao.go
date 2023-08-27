package threaddao

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/orpheus/strings/pkg/infra/sqldb"
)

type VersionedThreadDao struct {
	*sqldb.Store
}

// Save saves a new versioned thread record.
// Do not let the caller set archived, or deleted. These fields are handled by the service layer.
// Let date_created be set by the database.
func (t *VersionedThreadDao) Save(record *VersionedThreadRecord) (*VersionedThreadRecord, error) {
	query := `
	insert into versioned_thread (
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

// FindByThreadId finds the latest version of a thread by thread id. Returns nil if no record is found.
func (t *VersionedThreadDao) FindByThreadId(threadId uuid.UUID) (*VersionedThreadRecord, error) {
	// query to grab the latest version of the versioned thread
	// 1. create a sub-query to grab only the latest versions of the records
	// 2. then join the records to get all the fields back
	query := `
	select tr.* from versioned_thread tr
	join (
		select thread_id, max(version) as maxVersion
		from versioned_thread
		group by thread_id
	) latest on tr.thread_id = latest.thread_id and tr.version = latest.maxVersion
	where tr.thread_id = $1;
	`

	ex := t.Store.GetExecutor()
	row := ex.QueryRow(query, threadId)

	var r VersionedThreadRecord
	if err := row.Scan(&r.Id, &r.Name, &r.Version, &r.ThreadId, &r.Archived, &r.Deleted, &r.DateCreated); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan record: %s", err)
	}

	return &r, nil
}

func (t *VersionedThreadDao) FindAll() ([]*VersionedThreadRecord, error) {
	query := `
	select tr.* from versioned_thread tr
	join (
		select thread_id, max(version) as maxVersion
		from versioned_thread
		group by thread_id
	) latest on tr.thread_id = latest.thread_id and tr.version = latest.maxVersion
	`

	ex := t.Store.GetExecutor()
	rows, err := ex.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed query table: %s", err)
	}
	defer rows.Close()

	var records []*VersionedThreadRecord
	for rows.Next() {
		var r VersionedThreadRecord
		if err := rows.Scan(&r.Id, &r.Name, &r.Version, &r.ThreadId, &r.Archived, &r.Deleted, &r.DateCreated); err != nil {
			return nil, fmt.Errorf("failed to scan record: %s", err)
		}
		records = append(records, &r)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("err while iterating over rows: %s", err)
	}

	return records, nil
}
