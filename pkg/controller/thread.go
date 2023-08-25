package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/orpheus/strings/pkg/infra/sqldb"
	"github.com/orpheus/strings/pkg/repo/threads"
	"github.com/orpheus/strings/pkg/service"
	"net/http"
)

func NewThreadController(router *gin.RouterGroup, store *sqldb.Store) *ThreadController {
	controller := &ThreadController{
		ThreadService: service.NewThreadService(
			threads.NewThreadRepository(store),
		),
	}

	controller.RegisterRoutes(router)

	return controller
}

type ThreadController struct {
	ThreadService ThreadService
}

type ThreadService interface {
	PostThread(thread *threads.Thread) (*threads.Thread, error)
	GetThreads() ([]threads.Thread, error)
	GetThreadIds() ([]uuid.UUID, error) // used if ?only_ids=true
	ArchiveThread(id uuid.UUID) (*threads.Thread, error)
	RestoreThread(id uuid.UUID) (*threads.Thread, error)
	ActivateThread(id uuid.UUID) (*threads.Thread, error)
	DeactivateThread(id uuid.UUID) (*threads.Thread, error)
	DeleteThread(id uuid.UUID) (*threads.Thread, error)
}

func (t *ThreadController) RegisterRoutes(router *gin.RouterGroup) {
	threadsRouterGroup := router.Group("/threads")
	{
		threadsRouterGroup.POST("", t.PostThreads)
		threadsRouterGroup.GET("", t.GetThreads)
	}

	threadRouterGroup := router.Group("/thread")
	{
		threadRouterGroup.GET("", t.Archive)
		threadRouterGroup.POST("/archive/:id", t.Archive)
		threadRouterGroup.POST("/restore/:id", t.Restore)
		threadRouterGroup.POST("/activate/:id", t.Activate)
		threadRouterGroup.POST("/deactivate/:id", t.Deactivate)
		threadRouterGroup.POST("/delete/:id", t.Delete)
	}
}

func (t *ThreadController) PostThreads(c *gin.Context) {
	var thread threads.Thread

	// using c.BindJSON calls c.MustBindJSON which writes the response header as `text/plain` which can't be overriddens
	err := c.ShouldBindJSON(&thread)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewApiError(0, fmt.Sprintf("failed to bind request body with thread: %s", err)))
		return
	}

	threadsResponse, err := t.ThreadService.PostThread(&thread)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewApiError(0, fmt.Sprintf("ThreadService.PostThread failed: %s", err)))
	}

	c.JSON(http.StatusOK, threadsResponse)
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
