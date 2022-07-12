package core

import (
	"github.com/gofrs/uuid"
	"time"
)

type String struct {
	Id           uuid.UUID `json:"id"`
	Name         string    `json:"name" binding:"required"`
	Order        int       `json:"order"`
	Thread       uuid.UUID `json:"thread" binding:"required"`
	Description  string    `json:"description,omitempty"`
	DateCreated  time.Time `json:"dateCreated"`
	DateModified time.Time `json:"dateModified"`
}

type StringOrder struct {
	Id    uuid.UUID `json:"id"`
	Order int       `json:"order"`
}
