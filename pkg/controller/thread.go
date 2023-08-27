package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/orpheus/strings/pkg/core"
	"github.com/orpheus/strings/pkg/infra/sqldb"
	"github.com/orpheus/strings/pkg/persistence/dao/stringdao"
	"github.com/orpheus/strings/pkg/persistence/dao/threaddao"
	"github.com/orpheus/strings/pkg/persistence/repo/pgrepo/stringrepo"
	"github.com/orpheus/strings/pkg/persistence/repo/pgrepo/threadrepo"
	"github.com/orpheus/strings/pkg/service"
	"net/http"
)

func NewThreadController(router *gin.RouterGroup, store *sqldb.Store) *ThreadController {
	threadDao := &threaddao.ThreadDao{Store: store}
	versionedThreadDao := &threaddao.VersionedThreadDao{Store: store}

	stringDao := &stringdao.StringDao{Store: store}
	versionedStringDao := &stringdao.VersionedStringDao{Store: store}

	controller := &ThreadController{
		ThreadService: service.NewThreadService(
			threadrepo.NewThreadRepository(threadDao, versionedThreadDao),
			stringrepo.NewStringRepository(stringDao, versionedStringDao),
		),
	}

	controller.RegisterRoutes(router)

	return controller
}

type ThreadController struct {
	ThreadService ThreadService
}

type ThreadService interface {
	PostThread(thread *core.Thread) (*core.Thread, error)
	GetThreads() ([]*core.Thread, error)
	GetThreadIds() ([]uuid.UUID, error) // used if ?only_ids=true
	ArchiveThread(id uuid.UUID) (*core.Thread, error)
	RestoreThread(id uuid.UUID) (*core.Thread, error)
	ActivateThread(id uuid.UUID) (*core.Thread, error)
	DeactivateThread(id uuid.UUID) (*core.Thread, error)
	DeleteThread(id uuid.UUID) (*core.Thread, error)
}

func (t *ThreadController) RegisterRoutes(router *gin.RouterGroup) {
	threadsRouterGroup := router.Group("/threads")
	{
		threadsRouterGroup.POST("", t.PostThreads)
		threadsRouterGroup.GET("", t.GetThreads)
		threadsRouterGroup.POST("/archive/:id", t.Archive)
		threadsRouterGroup.POST("/restore/:id", t.Restore)
		threadsRouterGroup.POST("/activate/:id", t.Activate)
		threadsRouterGroup.POST("/deactivate/:id", t.Deactivate)
		threadsRouterGroup.POST("/delete/:id", t.Delete)
	}
}

func (t *ThreadController) PostThreads(c *gin.Context) {
	var thread core.Thread

	// using c.BindJSON calls c.MustBindJSON which writes the response header as `text/plain` which can't be overriddens
	err := c.ShouldBindJSON(&thread)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewApiError(0, fmt.Sprintf("failed to bind request body with thread: %s", err)))
		return
	}

	threadsResponse, err := t.ThreadService.PostThread(&thread)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewApiError(0, fmt.Sprintf("ThreadService.PostThread failed: %s", err)))
		return
	}

	c.JSON(http.StatusOK, threadsResponse)
}

func (t *ThreadController) GetThreads(c *gin.Context) {
	threads, err := t.ThreadService.GetThreads()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewApiError(0, fmt.Sprintf("ThreadService.GetThreads failed: %s", err)))
		return
	}

	c.JSON(http.StatusOK, threads)
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
