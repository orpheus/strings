package core

import (
	"github.com/gofrs/uuid"
	"time"
)

type Thread struct {
	Id           uuid.UUID `json:"id"`
	Name         string    `json:"name" binding:"required"`
	Description  string    `json:"description"`
	DateCreated  time.Time `json:"dateCreated"`
	DateModified time.Time `json:"dateModified"`
}
