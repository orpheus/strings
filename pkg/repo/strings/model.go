package strings

import (
	"github.com/gofrs/uuid"
	"time"
)

type StringId struct {
	Id uuid.UUID `json:"id"`
}

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
