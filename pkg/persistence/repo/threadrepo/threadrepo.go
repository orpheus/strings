package threadrepo

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/orpheus/strings/pkg/core"
	"github.com/orpheus/strings/pkg/persistence/dao/threaddao"
)

type Repository struct {
	ThreadDao
	StringDao
	VersionedThreadDao
}

type StringDao interface {
}

type ThreadDao interface {
	Save(record *threaddao.ThreadRecord) (*threaddao.ThreadRecord, error)
}

type VersionedThreadDao interface {
	Save(record *threaddao.VersionedThreadRecord) (*threaddao.VersionedThreadRecord, error)
	FindByThreadId(threadId uuid.UUID) (*threaddao.VersionedThreadRecord, error)
}

func NewThreadRepository(threadDao ThreadDao, stringDao StringDao, versionedThreadDao VersionedThreadDao) *Repository {
	return &Repository{
		ThreadDao:          threadDao,
		StringDao:          stringDao,
		VersionedThreadDao: versionedThreadDao,
	}
}

func (r *Repository) FindByThreadId(threadId uuid.UUID) (*core.Thread, error) {
	versionedThread, err := r.VersionedThreadDao.FindByThreadId(threadId)
	if err != nil {
		return nil, fmt.Errorf("failed to find by ThreadId: %s", err)
	}

	if versionedThread == nil {
		return nil, nil
	}

	// TODO("Get strings")

	return &core.Thread{
		Id:          versionedThread.Id,
		Name:        versionedThread.Name,
		Version:     versionedThread.Version,
		ThreadId:    versionedThread.ThreadId,
		Archived:    versionedThread.Archived,
		Deleted:     versionedThread.Deleted,
		DateCreated: versionedThread.DateCreated,
		Strings:     nil,
	}, nil
}

func (r *Repository) CreateThread(name string, id, threadId uuid.UUID) (*core.Thread, error) {
	if threadId == (uuid.UUID{}) {
		return nil, fmt.Errorf("missing ThreadId")
	}

	if id == (uuid.UUID{}) {
		return nil, fmt.Errorf("missing Id\n")
	}

	threadRecord, err := r.ThreadDao.Save(&threaddao.ThreadRecord{Id: threadId})
	if err != nil {
		return nil, fmt.Errorf("failed to create thread record: %s\n", err)
	}

	versionedThread := threaddao.VersionedThreadRecord{
		Id:       id,
		Name:     name,
		Version:  1,
		ThreadId: threadRecord.Id,
		// rest of the values will be defaulted by postgres
	}

	versionedThreadRecord, err := r.VersionedThreadDao.Save(&versionedThread)
	if err != nil {
		return nil, fmt.Errorf("failed to create versioned thread record: %s\n", err)
	}

	thread := &core.Thread{
		Id:          versionedThreadRecord.Id,
		Name:        versionedThreadRecord.Name,
		Version:     versionedThreadRecord.Version,
		ThreadId:    versionedThreadRecord.ThreadId,
		Archived:    versionedThreadRecord.Archived,
		Deleted:     versionedThreadRecord.Deleted,
		DateCreated: versionedThreadRecord.DateCreated,
		Strings:     nil,
	}

	return thread, nil
}
