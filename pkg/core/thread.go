package core

import (
	"github.com/google/uuid"
	"reflect"
	"time"
)

type Thread struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`    //  binding:"required"`
	Version     int       `json:"version"` //  binding:"required"`
	ThreadId    uuid.UUID `json:"thread_id"`
	Archived    bool      `json:"archived"`
	Deleted     bool      `json:"deleted"`
	DateCreated time.Time `json:"dateCreated"`
	Strings     []*String `json:"strings"`
}

func (t *Thread) Diff(other *Thread) bool {
	return t.DiffThreadOnly(other) || t.DiffStringsOnly(other)
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
func (t *Thread) DiffStringsOnly(other *Thread) bool {
	that := make(map[uuid.UUID]*String)

	for _, str := range other.Strings {
		that[str.Id] = str
	}

	for _, str := range t.Strings {
		if otherStr, exists := that[str.Id]; exists {
			if str.Diff(otherStr) {
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
