package errorhandler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/orpheus/strings/pkg/core"
	"github.com/orpheus/strings/pkg/persistence/repo/pgrepo/stringrepo"
	"github.com/orpheus/strings/pkg/persistence/repo/pgrepo/threadrepo"
	"github.com/orpheus/strings/pkg/service/threadsvc"
	"net/http"
)

type ApiError struct {
	ErrorDescriptor core.ErrorDescriptor `json:"errorDescriptor"`
}

func NewApiError(code int, msg string) ApiError {
	return ApiError{core.ErrorDescriptor{
		ErrorCode:    code,
		ErrorMessage: msg,
	}}
}

// Defined custom errors here that are the result of a client issue
var clientErrors = []error{
	threadrepo.ErrThreadAlreadyArchived,
	threadrepo.ErrThreadAlreadyDeleted,
	threadrepo.ErrThreadAlreadyRestored,

	stringrepo.ErrStringNotFound,
	stringrepo.ErrStringMissingName,
	stringrepo.ErrStringAlreadyDeleted,
	stringrepo.ErrStringAlreadyArchived,
	stringrepo.ErrStringAlreadyRestored,
	stringrepo.ErrStringAlreadyActive,
	stringrepo.ErrStringAlreadyDeleted,
	stringrepo.ErrStringAlreadyPrivate,
	stringrepo.ErrStringAlreadyPublic,

	threadsvc.ErrThreadCannotBeUpdated,
	threadsvc.ErrStringCannotBeUpdated,
}

// Defined custom errors here that are the result of a server issue
var serverErrors []error

func HandleApiError(c *gin.Context, err error) {
	for _, clientError := range clientErrors {
		if errors.Is(err, clientError) {
			c.AbortWithStatusJSON(http.StatusBadRequest, NewApiError(0, err.Error()))
			return
		}
	}

	for _, serverError := range serverErrors {
		if errors.Is(err, serverError) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, NewApiError(0, err.Error()))
			return
		}
	}

	c.AbortWithStatusJSON(http.StatusInternalServerError, NewApiError(0, err.Error()))
}
