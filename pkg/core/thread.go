package core

import (
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"time"
)

type Thread struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Version     int       `json:"version"`
	ThreadId    uuid.UUID `json:"thread_id"`
	Archived    bool      `json:"archived"`
	Deleted     bool      `json:"deleted"`
	DateCreated time.Time `json:"dateCreated"`
	Strings     []*String `json:"strings"`
}

// UpdateFromClientIgnoreStrings updates the thread with the values from the clientThread.
// Does not update strings.
func (t *Thread) UpdateFromClientIgnoreStrings(clientThread *Thread) {
	t.Name = clientThread.Name

	// handle strings separately because any new strings have to be created
	// and any existing strings have to be updated
}

func (t *Thread) ValidateSelf() error {
	if t.Name == "" {
		return fmt.Errorf("missing `name`")
	}

	return nil
}

func (t *Thread) Diff(other *Thread) bool {
	threadDif := t.DiffThreadOnly(other)
	stringsDif := t.DiffStringsOnly(other)

	return threadDif || stringsDif
}

// DiffThreadOnly returns true if there is a diff between thread-specific content.
// Does not compare strings.
func (t *Thread) DiffThreadOnly(other *Thread) bool {
	// Compare just thread-specific content

	this := &Thread{
		Name: t.Name,
	}

	that := &Thread{
		Name: other.Name,
	}

	return !reflect.DeepEqual(this, that)
}

// DiffStringsOnly returns true if there is a diff between string content.
func (t *Thread) DiffStringsOnly(otherThread *Thread) bool {
	otherThreadMap := make(map[uuid.UUID]*String)

	for _, str := range otherThread.Strings {
		otherThreadMap[str.StringId] = str
	}

	for _, thisString := range t.Strings {
		if otherString, exists := otherThreadMap[thisString.StringId]; exists {
			if thisString.Diff(otherString) {
				// return true if string in this thread is different than string in other thread
				return true
			}
		} else {
			// return true if this thread contains a string the other thread doesn't have
			return true
		}
	}

	// return false if no changes found between strings in this and that thread
	return false
}

func (t *Thread) Locked() bool {
	return t.Deleted == true || t.Archived == true
}
