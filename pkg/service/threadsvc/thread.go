package threadsvc

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/orpheus/strings/pkg/core"
)

func NewThreadService(threadRepository ThreadRepository, stringRepository StringRepository) *ThreadService {
	return &ThreadService{
		ThreadRepository: threadRepository,
		StringRepository: stringRepository,
	}
}

type ThreadService struct {
	ThreadRepository ThreadRepository
	StringRepository StringRepository
}

type ThreadRepository interface {
	FindByThreadId(threadId uuid.UUID) (*core.Thread, error)
	CreateNewThread(thread *core.Thread) (*core.Thread, error)
	CreateVersionedThread(thread *core.Thread) (*core.Thread, error)
	FindAll() ([]*core.Thread, error)
	DeleteByThreadId(threadId uuid.UUID) error
	ArchiveByThreadId(threadId uuid.UUID) error
	RestoreByThreadId(threadId uuid.UUID) error
}

type StringRepository interface {
	CreateNewString(string *core.String) (*core.String, error)
	CreateVersionedString(string *core.String) (*core.String, error)
	FindAllByThreadId(threadId uuid.UUID) ([]*core.String, error)
}

func (t *ThreadService) PostThread(thread *core.Thread) (*core.Thread, error) {
	if thread.ThreadId == uuid.Nil {
		return t.createNewThread(thread)
	}

	serverThread, err := t.getServerThread(thread.ThreadId)
	if err != nil {
		return nil, err
	}

	if serverThread == nil {
		return t.createNewThread(thread)
	} else {
		updatedThread, err := t.updateThreadIfNeeded(thread, serverThread)
		if err != nil {
			return nil, err
		}
		return updatedThread.FilterDeleted(), nil
	}
}

func (t *ThreadService) getServerThread(threadId uuid.UUID) (*core.Thread, error) {
	serverThread, err := t.ThreadRepository.FindByThreadId(threadId)
	if err != nil {
		return nil, fmt.Errorf("thread repository failed to find thread by ThreadId: %s", err)
	}

	if serverThread == nil {
		return nil, nil
	}

	serverStrings, err := t.StringRepository.FindAllByThreadId(threadId)
	if err != nil {
		return nil, fmt.Errorf("string repository failed to find strings by ThreadId: %s", err)
	}

	serverThread.Strings = serverStrings

	return serverThread, nil
}

// note: thread name+version combination should be unique
func (t *ThreadService) createNewThread(thread *core.Thread) (*core.Thread, error) {
	if thread.Name == "" {
		return nil, fmt.Errorf("failed to create new thread, missing `name`")
	}

	if thread.ThreadId == uuid.Nil {
		thread.ThreadId = uuid.New()
	}

	newThread, err := t.ThreadRepository.CreateNewThread(thread)
	if err != nil {
		return nil, fmt.Errorf("failed to create new thread: %s\n", err)
	}

	serverStrings, err := t.createStrings(thread)
	if err != nil {
		return nil, err
	}

	newThread.Strings = serverStrings

	return newThread, nil
}

var ErrThreadCannotBeUpdated = fmt.Errorf("deleted or archived threads cannot be updated")

func (t *ThreadService) updateThreadIfNeeded(clientThread *core.Thread, serverThread *core.Thread) (*core.Thread, error) {
	// if client did not provide name, just use the server name. client cannot set empty name
	if clientThread.Name == "" {
		clientThread.Name = serverThread.Name
	}

	// checks if there's a difference in the thread or strings
	if !clientThread.Diff(serverThread) {
		return serverThread, nil
	}

	// if they didn't try to update the thread, just return the server thread
	// but do not let the client try to update a deleted or archived thread
	if serverThread.Locked() {
		return nil, ErrThreadCannotBeUpdated
	}

	if serverThread == nil {
		return nil, fmt.Errorf("cannot update thread, thread not found %s", clientThread.Id)
	}

	// updates just the thread values (not strings)
	serverThread.UpdateFromClientIgnoreStrings(clientThread)

	// checks that the client thread updates are valid
	if err := serverThread.ValidateSelf(); err != nil {
		return nil, fmt.Errorf("thread failed validation, %s", err)
	}

	serverStrings, err := t.updateAndCreateStrings(clientThread, serverThread)
	if err != nil {
		return nil, err
	}

	updatedThread, err := t.ThreadRepository.CreateVersionedThread(serverThread)
	if err != nil {
		return nil, err
	}

	updatedThread.Strings = serverStrings

	return updatedThread, nil
}

func (t *ThreadService) createStrings(thread *core.Thread) ([]*core.String, error) {
	if thread.Strings == nil {
		return nil, nil
	}

	stringMap := NewStringMap(thread.Strings)
	err := stringMap.OrderStrings()
	if err != nil {
		return nil, err
	}

	var serverStrings []*core.String

	for _, stringItem := range stringMap.Default {
		stringItem.ThreadId = thread.ThreadId

		serverString, err := t.StringRepository.CreateNewString(stringItem)
		if err != nil {
			return nil, err
		}
		serverStrings = append(serverStrings, serverString)
	}

	return serverStrings, nil
}

var ErrStringCannotBeUpdated = fmt.Errorf("deleted or archived strings cannot be updated")

func (t *ThreadService) updateAndCreateStrings(clientThread, serverThread *core.Thread) ([]*core.String, error) {
	stringMap := NewStringMap(serverThread.Strings)
	updatedStrings, err := stringMap.UpdateFrom(clientThread.Strings)
	if err != nil {
		return nil, err
	}

	newStringsFromClient := stringMap.GetNewStrings(clientThread.Strings)

	stringMap.Include(newStringsFromClient)

	err = stringMap.OrderStrings()
	if err != nil {
		return nil, err
	}

	// Loop through ordered string and update existing strings and create new strings
	var serverStrings []*core.String

	for _, orderedString := range stringMap.Default {
		if _, exists := stringMap.sourceStrings[orderedString.StringId]; exists {
			if _, exists := updatedStrings[orderedString.StringId]; exists {
				updatedString, err := t.StringRepository.CreateVersionedString(orderedString)
				if err != nil {
					return nil, fmt.Errorf("failed to update string: %s", err)
				}
				serverStrings = append(serverStrings, updatedString)
			} else {
				serverStrings = append(serverStrings, orderedString)
			}

		} else {
			// set thread id
			orderedString.ThreadId = serverThread.ThreadId

			updatedString, err := t.StringRepository.CreateNewString(orderedString)
			if err != nil {
				return nil, err
			}
			serverStrings = append(serverStrings, updatedString)
		}
	}

	serverStrings = append(serverStrings, stringMap.Archived...)

	return serverStrings, nil
}

func (t *ThreadService) GetThreads() ([]*core.Thread, error) {
	threads, err := t.ThreadRepository.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get threads: %s", err)
	}

	// TODO: have ThreadRepository return threads with strings
	for _, thread := range threads {
		strings, err := t.StringRepository.FindAllByThreadId(thread.ThreadId)
		if err != nil {
			return nil, fmt.Errorf("failed to get strings for thread %s: %s", thread.ThreadId, err)
		}
		thread.Strings = strings
	}

	return threads, nil
}

func (t *ThreadService) GetThreadIds() ([]uuid.UUID, error) {
	//TODO implement me
	panic("implement me")
}

func (t *ThreadService) ArchiveThread(threadId uuid.UUID) error {
	return t.ThreadRepository.ArchiveByThreadId(threadId)
}

func (t *ThreadService) RestoreThread(threadId uuid.UUID) error {
	return t.ThreadRepository.RestoreByThreadId(threadId)
}

func (t *ThreadService) DeleteThread(threadId uuid.UUID) error {
	return t.ThreadRepository.DeleteByThreadId(threadId)
}
