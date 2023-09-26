package stringdao

import (
	"github.com/google/uuid"
	"github.com/orpheus/strings/pkg/core"
	"time"
)

type VersionedStringRecord struct {
	Id          uuid.UUID
	Name        string
	Version     int
	StringId    uuid.UUID
	ThreadId    uuid.UUID
	Order       int
	Active      bool
	Archived    bool
	Private     bool
	Deleted     bool
	DateCreated time.Time
}

func (v *VersionedStringRecord) ToString() *core.String {
	return &core.String{
		Id:          v.Id,
		Name:        v.Name,
		Version:     v.Version,
		StringId:    v.StringId,
		ThreadId:    v.ThreadId,
		Order:       v.Order,
		Active:      v.Active,
		Archived:    v.Archived,
		Private:     v.Private,
		Deleted:     v.Deleted,
		DateCreated: v.DateCreated,
	}
}
