package stringdao

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/orpheus/strings/pkg/infra/sqldb"
)

type VersionedStringDao struct {
	*sqldb.Store
}

// Save saves a new versioned string record.
// Do not let the caller set active, archived, or deleted. These fields are handled by the service layer.
// Let date_created be set by the database.
func (v *VersionedStringDao) Save(record *VersionedStringRecord) (*VersionedStringRecord, error) {
	if record.ThreadId == uuid.Nil {
		return nil, errors.New("thread id cannot be nil")
	}

	if record.StringId == uuid.Nil {
		return nil, errors.New("string id cannot be nil")
	}

	query := `
	insert into versioned_string (
		id, name, version, string_id, thread_id, "order", active, archived, private, deleted
	) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
	returning id, name, version, string_id, thread_id, "order", active, archived, private, deleted, date_created;
	`

	ex := v.Store.GetExecutor()
	row := ex.QueryRow(query, record.Id, record.Name, record.Version, record.StringId, record.ThreadId, record.Order, record.Active, record.Archived, record.Private, record.Deleted)

	return scanRecord(row)
}

// FindByStringId finds the latest version of a string by string id. Returns nil if no record is found.
func (v *VersionedStringDao) FindByStringId(stringId uuid.UUID) (*VersionedStringRecord, error) {
	query := `
	select str.* from versioned_string str
	join (
		select string_id, max(version) as maxVersion
		from versioned_string
		group by string_id
	) latest on str.string_id = latest.string_id and str.version = latest.maxVersion
	where str.string_id = $1;
	`

	ex := v.Store.GetExecutor()
	row := ex.QueryRow(query, stringId)

	return scanRecord(row)
}

func (v *VersionedStringDao) FindAllByThreadId(threadId uuid.UUID) ([]*VersionedStringRecord, error) {
	if threadId == uuid.Nil {
		return nil, errors.New("thread id nil")
	}

	query := `
	select str.* from versioned_string str
	join (
		select string_id, max(version) as maxVersion
		from versioned_string
		group by string_id
	) latest on str.string_id = latest.string_id and str.version = latest.maxVersion
	where str.thread_id = $1;
	`

	ex := v.Store.GetExecutor()
	rows, err := ex.Query(query, threadId)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %s", err)
	}

	return scanRecords(rows)
}

func (v *VersionedStringDao) DeprecatedFindAllInThreadByStringId(stringId uuid.UUID) ([]*VersionedStringRecord, error) {
	if stringId == uuid.Nil {
		return nil, errors.New("stringId id nil")
	}

	query := `
	select str.*
	from versioned_string str
			 join(select string_id, max(version) as maxVersion
				  from versioned_string
				  group by string_id) latest
				 on str.string_id = latest.string_id and str.version = latest.maxVersion
	where str.thread_id = (select thread_id
						   from versioned_string
						   where versioned_string.string_id = $1
						   order by version desc
						   limit 1)
	order by "order" asc
	`

	ex := v.Store.GetExecutor()
	rows, err := ex.Query(query, stringId)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %s", err)
	}

	return scanRecords(rows)
}

func scanMe() (record *VersionedStringRecord, pointsToFields []any) {
	var r VersionedStringRecord
	return &r, []any{
		&r.Id,
		&r.Name,
		&r.Version,
		&r.StringId,
		&r.ThreadId,
		&r.Order,
		&r.Active,
		&r.Archived,
		&r.Private,
		&r.Deleted,
		&r.DateCreated,
	}
}

func scanRecords(rows *sql.Rows) ([]*VersionedStringRecord, error) {
	var records []*VersionedStringRecord
	for rows.Next() {
		r, fields := scanMe()
		if err := rows.Scan(fields...); err != nil {
			return nil, fmt.Errorf("failed to scan record: %s", err)
		}

		records = append(records, r)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows err: %s", rows.Err())
	}

	return records, nil
}

func scanRecord(row *sql.Row) (*VersionedStringRecord, error) {
	r, fields := scanMe()
	err := row.Scan(fields...)
	if err != nil {
		return nil, fmt.Errorf("scan err: %s", err)
	}

	return r, nil
}
