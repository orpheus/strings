package stringdao

import (
	"github.com/google/uuid"
	"github.com/orpheus/strings/pkg/infra/sqldb"
	"time"
)

type VersionedStringRecord struct {
	Id          uuid.UUID
	Version     int
	StringId    uuid.UUID
	ThreadId    uuid.UUID
	Order       int
	Active      bool
	Archived    bool
	Deleted     bool
	DateCreated time.Time
}

type VersionedStringDao struct {
	*sqldb.Store
}
