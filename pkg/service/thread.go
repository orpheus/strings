package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/orpheus/strings/pkg/core"
	"sort"
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
	CreateNewThreadVersion(thread *core.Thread) (*core.Thread, error)
	FindAll() ([]*core.Thread, error)
}

type StringRepository interface {
	CreateNewString(string *core.String) (*core.String, error)
	CreateNewStringVersion(string *core.String) (*core.String, error)
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
		return t.updateThreadIfNeeded(thread, serverThread)
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

func (t *ThreadService) updateThreadIfNeeded(clientThread *core.Thread, serverThread *core.Thread) (*core.Thread, error) {
	// if client did not provide name, just use the server name. client cannot set empty name
	if clientThread.Name == "" {
		clientThread.Name = serverThread.Name
	}

	if !clientThread.Diff(serverThread) {
		return serverThread, nil
	}

	if serverThread == nil {
		return nil, fmt.Errorf("cannot update thread, thread not found %s", clientThread.Id)
	}

	serverThread.UpdateFromClientIgnoreStrings(clientThread)

	if err := serverThread.ValidateSelf(); err != nil {
		return nil, fmt.Errorf("thread failed validation, %s", err)
	}

	serverStrings, err := t.updateAndCreateStrings(clientThread, serverThread)
	if err != nil {
		return nil, err
	}

	serverThread.Strings = serverStrings

	return t.ThreadRepository.CreateNewThreadVersion(serverThread)
}

func (t *ThreadService) createStrings(thread *core.Thread) ([]*core.String, error) {
	if thread.Strings == nil {
		return nil, nil
	}

	var orderedStrings []*core.String
	var unorderedStrings []*core.String

	for _, stringItem := range thread.Strings {
		if stringItem.Order != 0 {
			orderedStrings = append(orderedStrings, stringItem)
		} else {
			unorderedStrings = append(unorderedStrings, stringItem)
		}
	}

	// Sort ordered strings by order

	sort.Slice(orderedStrings, func(i, j int) bool {
		return orderedStrings[i].Order < orderedStrings[j].Order
	})

	// Order and validate order

	maxRange := len(orderedStrings) + len(unorderedStrings)

	for index, stringItem := range orderedStrings {
		if stringItem.Order > maxRange {
			return nil, fmt.Errorf("string (%s) order is greater than max range %d", stringItem.StringId, maxRange)
		}

		if stringItem.Order != index+1 {
			return nil, fmt.Errorf("invalid string order")
		}
	}

	for _, stringItem := range unorderedStrings {
		stringItem.Order = len(orderedStrings) + 1
		orderedStrings = append(orderedStrings, stringItem)
	}

	var serverStrings []*core.String

	for _, stringItem := range orderedStrings {
		stringItem.ThreadId = thread.ThreadId

		serverString, err := t.StringRepository.CreateNewString(stringItem)
		if err != nil {
			return nil, fmt.Errorf("failed to create string: %s", err)
		}
		serverStrings = append(serverStrings, serverString)
	}

	return serverStrings, nil
}

func (t *ThreadService) updateAndCreateStrings(clientThread, serverThread *core.Thread) ([]*core.String, error) {
	serverStringMap := make(map[uuid.UUID]*core.String)

	for _, serverString := range serverThread.Strings {
		serverStringMap[serverString.StringId] = serverString
	}

	// map of strings that have been updated
	updatedStrings := make(map[uuid.UUID]*core.String)

	// Update-mutate existing strings and mark as updated only if they have changed

	for _, clientString := range clientThread.Strings {
		if _, exists := serverStringMap[clientString.StringId]; exists {
			serverString := serverStringMap[clientString.StringId]

			if _, exists := updatedStrings[serverString.StringId]; exists {
				return nil, fmt.Errorf("duplicate client string provided for string id %s", serverString.StringId)
			}

			if serverString.Diff(clientString) {
				serverString.UpdateFromClient(clientString)
				updatedStrings[serverString.StringId] = serverString
			}
		}
	}

	var newStrings []*core.String

	// Filter new strings

	for _, clientString := range clientThread.Strings {
		if _, exists := serverStringMap[clientString.StringId]; !exists {
			newStrings = append(newStrings, clientString)
		}
	}

	// Order and validate order

	maxRange := len(serverStringMap) + len(newStrings)

	var orderedStrings []*core.String
	var newUnorderedStrings []*core.String

	for _, serverString := range serverStringMap {
		if serverString.Order > maxRange {
			return nil, fmt.Errorf("string (%s) order is greater than max range %d", serverString.StringId, maxRange)
		}
		orderedStrings = append(orderedStrings, serverString)
	}

	for _, newString := range newStrings {
		if newString.Order != 0 {
			if newString.Order > maxRange {
				return nil, fmt.Errorf("string (%s) order is greater than max range %d", newString.StringId, maxRange)
			}
			orderedStrings = append(orderedStrings, newString)
		} else {
			newUnorderedStrings = append(newUnorderedStrings, newString)
		}
	}

	// Sort ordered strings by order

	sort.Slice(orderedStrings, func(i, j int) bool {
		return orderedStrings[i].Order < orderedStrings[j].Order
	})

	// Validate that the known order supplied is sequential

	for index, stringItem := range orderedStrings {
		if stringItem.Order != index+1 {
			return nil, fmt.Errorf("invalid string order")
		}
	}

	// Assign order to new strings and add them to ordered strings
	for _, newString := range newUnorderedStrings {
		newString.Order = len(orderedStrings) + 1
		orderedStrings = append(orderedStrings, newString)
	}

	// Loop through ordered string and update existing strings and create new strings
	var updatedServerStrings []*core.String

	for _, orderedString := range orderedStrings {
		if _, exists := serverStringMap[orderedString.StringId]; exists {
			if _, exists := updatedStrings[orderedString.StringId]; exists {
				updatedString, err := t.StringRepository.CreateNewStringVersion(orderedString)
				if err != nil {
					return nil, fmt.Errorf("failed to update string: %s", err)
				}
				updatedServerStrings = append(updatedServerStrings, updatedString)
			} else {
				updatedServerStrings = append(updatedServerStrings, orderedString)
			}

		} else {
			// set thread id
			orderedString.ThreadId = serverThread.ThreadId

			updatedString, err := t.StringRepository.CreateNewString(orderedString)
			if err != nil {
				return nil, fmt.Errorf("failed to create string: %s", err)
			}
			updatedServerStrings = append(updatedServerStrings, updatedString)
		}
	}

	return updatedServerStrings, nil
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
