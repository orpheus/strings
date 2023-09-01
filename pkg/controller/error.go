package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/orpheus/strings/pkg/core"
	"github.com/orpheus/strings/pkg/persistence/repo/pgrepo/threadrepo"
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

func handleServerHandlerError(c *gin.Context, err error) {
	clientErrors := []error{
		threadrepo.ErrThreadAlreadyArchived,
		threadrepo.ErrThreadAlreadyDeleted,
		threadrepo.ErrThreadAlreadyRestored,
	}

	var serverErrors []error

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
