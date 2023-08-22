package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/orpheus/strings/core"
	"github.com/orpheus/strings/infrastructure/logging"
	"log"
	"net/http"
)

type Controller struct {
	Interactor Interactor
	Logger     logging.Logger
}

type Interactor interface {
	FindAll() ([]core.Thread, error)
	CreateOne(thread core.Thread) (core.Thread, error)
	DeleteById(id uuid.UUID) error
}

func (s *Controller) RegisterRoutes(router *gin.RouterGroup) {
	thread := router.Group("/thread")
	{
		thread.GET("", s.FindAll)
		thread.POST("", s.CreateOne)
		thread.DELETE("/:id", s.DeleteById)
	}
}

func (s *Controller) FindAll(c *gin.Context) {
	threads, err := s.Interactor.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, threads)
}

func (s *Controller) CreateOne(c *gin.Context) {
	var thread core.Thread
	if err := c.ShouldBindJSON(&thread); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	newThread, err := s.Interactor.CreateOne(thread)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("%s", err.Error())},
		)
		return
	}
	c.IndentedJSON(http.StatusOK, newThread)
}

func (s *Controller) DeleteById(c *gin.Context) {
	threadId, err := uuid.FromString(c.Param("id"))
	if err != nil {
		log.Fatalf("failed to parse UUID %q: %v", s, err)
	}

	err = s.Interactor.DeleteById(threadId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, true)
}
