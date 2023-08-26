package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/orpheus/strings/pkg/core"
)

func NewThreadService(threadRepository ThreadRepository) *ThreadService {
	return &ThreadService{
		ThreadRepository: threadRepository,
	}
}

type ThreadService struct {
	ThreadRepository ThreadRepository
}

type ThreadRepository interface {
	FindByThreadId(threadId uuid.UUID) (*core.Thread, error)
	CreateThread(name string, id, threadId uuid.UUID) (*core.Thread, error)
}

func (t *ThreadService) PostThread(thread *core.Thread) (*core.Thread, error) {
	if thread.Id == (uuid.UUID{}) {
		return t.createNewThread(thread)
	}

	existingThread, err := t.ThreadRepository.FindByThreadId(thread.ThreadId)
	if err != nil {
		return nil, fmt.Errorf("failed to find thread by threadId: %s\n", err)
	}

	if existingThread == nil {
		return t.createNewThread(thread)
	}

	return existingThread, nil
}

// note: thread name+version combination should be unique
func (t *ThreadService) createNewThread(thread *core.Thread) (*core.Thread, error) {
	if thread.Name == "" {
		return nil, fmt.Errorf("failed to create new thread, missing `name`")
	}

	if thread.Id == (uuid.UUID{}) {
		thread.Id = uuid.New()
	}

	if thread.ThreadId == (uuid.UUID{}) {
		thread.ThreadId = uuid.New()
	}

	newThread, err := t.ThreadRepository.CreateThread(thread.Name, thread.Id, thread.ThreadId)
	if err != nil {
		return nil, fmt.Errorf("failed to create new thread: %s\n", err)
	}

	return newThread, nil
}

func (t *ThreadService) GetThreads() ([]core.Thread, error) {
	//TODO implement me
	panic("implement me")
}

func (t *ThreadService) GetThreadIds() ([]uuid.UUID, error) {
	//TODO implement me
	panic("implement me")
}

func (t *ThreadService) ArchiveThread(id uuid.UUID) (*core.Thread, error) {
	//TODO implement me
	panic("implement me")
}

func (t *ThreadService) RestoreThread(id uuid.UUID) (*core.Thread, error) {
	//TODO implement me
	panic("implement me")
}

func (t *ThreadService) ActivateThread(id uuid.UUID) (*core.Thread, error) {
	//TODO implement me
	panic("implement me")
}

func (t *ThreadService) DeactivateThread(id uuid.UUID) (*core.Thread, error) {
	//TODO implement me
	panic("implement me")
}

func (t *ThreadService) DeleteThread(id uuid.UUID) (*core.Thread, error) {
	//TODO implement me
	panic("implement me")
}
