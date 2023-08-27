package threaddao

import (
	"github.com/google/uuid"
	"github.com/orpheus/strings/pkg/core"
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

// FromThread converts a Thread to a VersionedThreadRecord
func (v *VersionedThreadRecord) FromThread(thread *core.Thread) *VersionedThreadRecord {
	return &VersionedThreadRecord{
		Id:          thread.Id,
		Name:        thread.Name,
		Version:     thread.Version,
		ThreadId:    thread.Id,
		Archived:    thread.Archived,
		Deleted:     thread.Deleted,
		DateCreated: thread.DateCreated,
	}
}

func (v *VersionedThreadRecord) ToThread(strings []*core.String) *core.Thread {
	return &core.Thread{
		Id:          v.Id,
		Name:        v.Name,
		Version:     v.Version,
		ThreadId:    v.ThreadId,
		Archived:    v.Archived,
		Deleted:     v.Deleted,
		DateCreated: v.DateCreated,
		Strings:     strings,
	}
}
