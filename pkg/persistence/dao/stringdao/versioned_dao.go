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
	query := `
	insert into versioned_string (
		id, name, version, string_id, thread_id, "order"
	) values ($1, $2, $3, $4) 
	returning id, name, version, string_id, thread_id, "order", active, archived, deleted, date_created;
	`

	ex := v.Store.GetExecutor()
	row := ex.QueryRow(query, record.Id, record.Name, record.Version, record.StringId, record.ThreadId, record.Order)

	var r VersionedStringRecord
	err := row.Scan(&r.Id, &r.Name, &r.Version, &r.StringId, &r.ThreadId, &r.Order, &r.Active, &r.Archived, &r.Deleted, &r.DateCreated)
	if err != nil {
		return nil, fmt.Errorf("scan err for versioned string record: %s", err)
	}

	return &r, nil
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

	var r VersionedStringRecord
	if err := row.Scan(&r.Id, &r.Name, &r.Version, &r.StringId, &r.ThreadId, &r.Order, &r.Active, &r.Archived, &r.Deleted, &r.DateCreated); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan record: %s", err)
	}

	return &r, nil
}
