package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/orpheus/strings/pkg/repo/threads"
	"time"
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
	FindByThreadId(threadId uuid.UUID) (*threads.Thread, error)
	CreateThread(thread *threads.Thread) (*threads.Thread, error)
}

func (t *ThreadService) PostThread(thread *threads.Thread) (*threads.Thread, error) {
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

func CreateNewThread(name string) *threads.Thread {
	return &threads.Thread{
		Id:          uuid.New(),
		Name:        name,
		Version:     1,
		ThreadId:    uuid.New(),
		Archived:    false,
		Deleted:     false,
		DateCreated: time.Now(),
		Strings:     nil,
	}
}

// note: thread name+version combination should be unique
func (t *ThreadService) createNewThread(thread *threads.Thread) (*threads.Thread, error) {
	if thread.Name == "" {
		return nil, fmt.Errorf("failed to create new thread, missing `name`")
	}

	newThread, err := t.ThreadRepository.CreateThread(CreateNewThread(thread.Name))
	if err != nil {
		return nil, fmt.Errorf("failed to create new thread: %s\n", err)
	}

	return newThread, nil
}

func (t *ThreadService) GetThreads() ([]threads.Thread, error) {
	//TODO implement me
	panic("implement me")
}

func (t *ThreadService) GetThreadIds() ([]uuid.UUID, error) {
	//TODO implement me
	panic("implement me")
}

func (t *ThreadService) ArchiveThread(id uuid.UUID) (*threads.Thread, error) {
	//TODO implement me
	panic("implement me")
}

func (t *ThreadService) RestoreThread(id uuid.UUID) (*threads.Thread, error) {
	//TODO implement me
	panic("implement me")
}

func (t *ThreadService) ActivateThread(id uuid.UUID) (*threads.Thread, error) {
	//TODO implement me
	panic("implement me")
}

func (t *ThreadService) DeactivateThread(id uuid.UUID) (*threads.Thread, error) {
	//TODO implement me
	panic("implement me")
}

func (t *ThreadService) DeleteThread(id uuid.UUID) (*threads.Thread, error) {
	//TODO implement me
	panic("implement me")
}
