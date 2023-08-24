package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/orpheus/strings/pkg/repo/threads"
)

type ThreadController struct {
	ThreadService ThreadService
}

type ThreadService interface {
	PostThread(thread threads.Thread) (threads.Thread, error)
	GetThreads() ([]threads.Thread, error)
	GetThreadIds() ([]uuid.UUID, error) // used if ?only_ids=true
	ArchiveThread(id uuid.UUID) (threads.Thread, error)
	RestoreThread(id uuid.UUID) (threads.Thread, error)
	ActivateThread(id uuid.UUID) (threads.Thread, error)
	DeactivateThread(id uuid.UUID) (threads.Thread, error)
	DeleteThread(id uuid.UUID) (threads.Thread, error)
}

func (t *ThreadController) RegisterRoutes(router *gin.RouterGroup) {
	threads := router.Group("/threads")
	{
		threads.POST("", t.PostThreads)
		threads.GET("", t.GetThreads)
	}

	thread := router.Group("/thread")
	{
		thread.GET("", t.Archive)
		thread.POST("/archive/:id", t.Archive)
		thread.POST("/restore/:id", t.Restore)
		thread.POST("/activate/:id", t.Activate)
		thread.POST("/deactivate/:id", t.Deactivate)
		thread.POST("/delete/:id", t.Delete)
	}
}

func (t *ThreadController) PostThreads(c *gin.Context) {

}

func (t *ThreadController) GetThreads(c *gin.Context) {

}

func (t *ThreadController) Archive(c *gin.Context) {

}

func (t *ThreadController) Restore(c *gin.Context) {

}

func (t *ThreadController) Activate(c *gin.Context) {

}

func (t *ThreadController) Deactivate(c *gin.Context) {

}

func (t *ThreadController) Delete(c *gin.Context) {

}
