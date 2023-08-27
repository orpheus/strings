package threadrepo

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/orpheus/strings/pkg/core"
	"github.com/orpheus/strings/pkg/persistence/dao/threaddao"
)

type Repository struct {
	ThreadDao
	VersionedThreadDao
}

type ThreadDao interface {
	Save(record *threaddao.ThreadRecord) (*threaddao.ThreadRecord, error)
}

type VersionedThreadDao interface {
	Save(record *threaddao.VersionedThreadRecord) (*threaddao.VersionedThreadRecord, error)
	FindByThreadId(threadId uuid.UUID) (*threaddao.VersionedThreadRecord, error)
	FindAll() ([]*threaddao.VersionedThreadRecord, error)
}

// golang: build exactly what you need and put them together.
// clojure: still the best language right now. any lisp is greater than a non-self-generating language.
// lisp allows you to self generate key mechanisms and build data pipes. you build your own system in which you can program in.
// any way you can think you can write something into the language itself and use it to further program your own program.

func NewThreadRepository(threadDao ThreadDao, versionedThreadDao VersionedThreadDao) *Repository {
	return &Repository{
		ThreadDao:          threadDao,
		VersionedThreadDao: versionedThreadDao,
	}
}

func (r *Repository) FindAll() ([]*core.Thread, error) {
	versionedThreads, err := r.VersionedThreadDao.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to find threads from versioned_thread table: %s", err)
	}

	// should the Thread Repository also fetch strings? or should that be done outside of this scope?

	return convertVersionedThreadsToCoreThreads(versionedThreads), nil
}

func (r *Repository) FindByThreadId(threadId uuid.UUID) (*core.Thread, error) {
	versionedThread, err := r.VersionedThreadDao.FindByThreadId(threadId)
	if err != nil {
		return nil, fmt.Errorf("failed to find by ThreadId: %s", err)
	}

	if versionedThread == nil {
		return nil, nil
	}

	// do not get strings here, that is a separate call

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

func (r *Repository) CreateNewThread(thread *core.Thread) (*core.Thread, error) {
	if thread == nil {
		return nil, fmt.Errorf("thread is nil")
	}

	if thread.ThreadId == uuid.Nil {
		thread.ThreadId = uuid.New()
	}

	if thread.Name == "" {
		return nil, fmt.Errorf("missing `name`")
	}

	threadRecord, err := r.ThreadDao.Save(&threaddao.ThreadRecord{Id: thread.ThreadId})
	if err != nil {
		return nil, fmt.Errorf("failed to create thread record: %s\n", err)
	}

	versionedThread := threaddao.VersionedThreadRecord{
		Id:       uuid.New(),
		Name:     thread.Name,
		Version:  1,
		ThreadId: threadRecord.Id,
		// rest of the values will be defaulted by postgres
	}

	versionedThreadRecord, err := r.VersionedThreadDao.Save(&versionedThread)
	if err != nil {
		return nil, fmt.Errorf("failed to create versioned thread record: %s\n", err)
	}

	return versionedThreadRecord.ToThread(thread.Strings), nil
}

// CreateNewThreadVersion updates an existing thread, does not care about or deal with any initial creation
// logic. That should be handled outside and separate from this function. Updates existing thread.
func (r *Repository) CreateNewThreadVersion(thread *core.Thread) (*core.Thread, error) {
	if thread == nil {
		return nil, fmt.Errorf("thread is nil")
	}

	if thread.ThreadId == uuid.Nil {
		return nil, fmt.Errorf("missing thread id")
	}

	serverThread, err := r.FindByThreadId(thread.ThreadId)
	if err != nil {
		return nil, fmt.Errorf("err during thread lookup: %s", err)
	}

	if serverThread == nil {
		return nil, fmt.Errorf("thread not found %s", thread.Id)
	}

	savedVersionedThreadRecord, err := r.VersionedThreadDao.Save(&threaddao.VersionedThreadRecord{
		Id:          uuid.New(),
		Name:        thread.Name,
		Version:     serverThread.Version + 1,
		ThreadId:    serverThread.ThreadId,
		Archived:    serverThread.Archived,
		Deleted:     serverThread.Deleted,
		DateCreated: serverThread.DateCreated,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new thread version: %s", err)
	}

	return savedVersionedThreadRecord.ToThread(thread.Strings), nil
}

func convertVersionedThreadsToCoreThreads(versionedThreads []*threaddao.VersionedThreadRecord) []*core.Thread {
	var threads []*core.Thread

	for _, t := range versionedThreads {
		// TODO: Need strings? right now the ThreadService manages getting and setting the strings
		threads = append(threads, t.ToThread(nil))
	}

	return threads
}
