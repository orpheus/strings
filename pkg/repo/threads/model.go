package threads

import (
	"github.com/google/uuid"
	"github.com/orpheus/strings/pkg/repo/strings"
	"time"
)

type ThreadId struct {
	Id uuid.UUID `json:"id"`
}

type Thread struct {
	Id          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`    //  binding:"required"`
	Version     int              `json:"version"` //  binding:"required"`
	ThreadId    uuid.UUID        `json:"thread_id"`
	Archived    bool             `json:"archived"`
	Deleted     bool             `json:"deleted"`
	DateCreated time.Time        `json:"dateCreated"`
	Strings     []strings.String `json:"strings"`
}
