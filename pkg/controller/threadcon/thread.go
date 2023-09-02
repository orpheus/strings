package threadcon

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/orpheus/strings/pkg/controller/errorhandler"
	"github.com/orpheus/strings/pkg/core"
	"github.com/orpheus/strings/pkg/infra/sqldb"
	"github.com/orpheus/strings/pkg/persistence/dao/stringdao"
	"github.com/orpheus/strings/pkg/persistence/dao/threaddao"
	"github.com/orpheus/strings/pkg/persistence/repo/pgrepo/stringrepo"
	"github.com/orpheus/strings/pkg/persistence/repo/pgrepo/threadrepo"
	"github.com/orpheus/strings/pkg/service/threadsvc"
	"net/http"
)

func NewThreadController(router *gin.RouterGroup, store *sqldb.Store) *ThreadController {
	threadDao := &threaddao.ThreadDao{Store: store}
	versionedThreadDao := &threaddao.VersionedThreadDao{Store: store}

	stringDao := &stringdao.StringDao{Store: store}
	versionedStringDao := &stringdao.VersionedStringDao{Store: store}

	controller := &ThreadController{
		ThreadService: threadsvc.NewThreadService(
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
	DeleteThread(threadId uuid.UUID) error
	ArchiveThread(threadId uuid.UUID) error
	RestoreThread(threadId uuid.UUID) error
}

func (t *ThreadController) RegisterRoutes(router *gin.RouterGroup) {
	threadsRouterGroup := router.Group("/threads")
	{
		threadsRouterGroup.POST("", t.PostThreads)
		threadsRouterGroup.GET("", t.GetThreads)
		threadsRouterGroup.POST("/archive/:id", t.Archive)
		threadsRouterGroup.POST("/restore/:id", t.Restore)
		threadsRouterGroup.POST("/delete/:id", t.Delete)
	}
}

func (t *ThreadController) PostThreads(c *gin.Context) {
	var thread core.Thread

	// using c.BindJSON calls c.MustBindJSON which writes the response header as `text/plain` which can't be overriddens
	err := c.ShouldBindJSON(&thread)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorhandler.NewApiError(0, fmt.Sprintf("failed to bind request body with thread: %s", err)))
		return
	}

	threadsResponse, err := t.ThreadService.PostThread(&thread)
	if err != nil {
		errorhandler.HandleApiError(c, err)
		return
	}

	if threadsResponse.Strings == nil {
		threadsResponse.Strings = []*core.String{}
	}

	c.JSON(http.StatusOK, threadsResponse)
}

func (t *ThreadController) GetThreads(c *gin.Context) {
	threads, err := t.ThreadService.GetThreads()
	if err != nil {
		errorhandler.HandleApiError(c, err)
		return
	}

	// TODO: optimize by having its own sql query
	onlyIds := c.Query("only_ids")
	if onlyIds == "true" {
		var ids []uuid.UUID
		for _, thread := range threads {
			ids = append(ids, thread.ThreadId)
		}
		c.JSON(http.StatusOK, ids)
		return
	}

	if threads == nil {
		threads = []*core.Thread{}
	}

	c.JSON(http.StatusOK, threads)
}

func (t *ThreadController) Archive(c *gin.Context) {
	threadId := c.Param("id")
	threadIdUuid, err := uuid.Parse(threadId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorhandler.NewApiError(0, fmt.Sprintf("%s", err)))
		return
	}

	err = t.ThreadService.ArchiveThread(threadIdUuid)
	if err != nil {
		errorhandler.HandleApiError(c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (t *ThreadController) Restore(c *gin.Context) {
	threadId := c.Param("id")
	threadIdUuid, err := uuid.Parse(threadId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorhandler.NewApiError(0, fmt.Sprintf("%s", err)))
		return
	}

	err = t.ThreadService.RestoreThread(threadIdUuid)
	if err != nil {
		errorhandler.HandleApiError(c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (t *ThreadController) Delete(c *gin.Context) {
	threadId := c.Param("id")
	threadIdUuid, err := uuid.Parse(threadId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorhandler.NewApiError(0, fmt.Sprintf("%s", err)))
		return
	}

	err = t.ThreadService.DeleteThread(threadIdUuid)
	if err != nil {
		errorhandler.HandleApiError(c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
