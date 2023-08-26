package core

import (
	"github.com/google/uuid"
	"reflect"
	"time"
)

type String struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`    //  binding:"required"`
	Version     int       `json:"version"` //  binding:"required"`
	StringId    uuid.UUID `json:"string_id"`
	ThreadId    uuid.UUID `json:"thread_id"`
	Order       int       `json:"order"`
	Active      bool      `json:"active"`
	Archived    bool      `json:"archived"`
	Deleted     bool      `json:"deleted"`
	DateCreated time.Time `json:"dateCreated"`
}

// Diff returns true if the string content does not match
func (s *String) Diff(other *String) bool {
	this := &String{
		Name:  s.Name,
		Order: s.Order,
	}

	that := &String{
		Name:  other.Name,
		Order: other.Order,
	}

	return !reflect.DeepEqual(this, that)
}
