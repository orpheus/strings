package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/orpheus/strings/pkg/repo/threads"
)

type StringController struct {
	StringsService StringsService
}

type StringsService interface {
	ArchiveString(id uuid.UUID) (threads.Thread, error)
	RestoreString(id uuid.UUID) (threads.Thread, error)
	ActivateString(id uuid.UUID) (threads.Thread, error)
	DeactivateString(id uuid.UUID) (threads.Thread, error)
	DeleteString(id uuid.UUID) (threads.Thread, error)
}

func (s *StringController) RegisterRoutes(router *gin.RouterGroup) {
	strings := router.Group("/strings")
	{
		strings.POST("/archive/:id", s.Archive)
		strings.POST("/restore/:id", s.Restore)
		strings.POST("/activate/:id", s.Activate)
		strings.POST("/deactivate/:id", s.Deactivate)
		strings.POST("/delete/:id", s.Delete)
	}
}

func (s *StringController) Archive(c *gin.Context) {

}

func (s *StringController) Restore(c *gin.Context) {

}

func (s *StringController) Activate(c *gin.Context) {

}

func (s *StringController) Deactivate(c *gin.Context) {

}

func (s *StringController) Delete(c *gin.Context) {

}
