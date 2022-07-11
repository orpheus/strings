package string

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/orpheus/strings/core"
	"github.com/orpheus/strings/infrastructure/logging"
	"log"
	"net/http"
)

// StringController is the api object that allows for http requests
// into our application
type StringController struct {
	Interactor StringInteractor
	Logger     logging.Logger
}

// StringInteractor defines the service interface the controller will usee
type StringInteractor interface {
	FindAll() ([]core.String, error)
	FindAllByThread(threadId uuid.UUID) ([]core.String, error)
	CreateOne(core.String) (core.String, error)
	DeleteById(id uuid.UUID) error
}

// RegisterRoutes creates a gin route grouping for the `/string` routes
func (s *StringController) RegisterRoutes(router *gin.RouterGroup) {
	skill := router.Group("/string")
	{
		skill.GET("", s.FindAll)
		skill.POST("", s.CreateOne)
		skill.DELETE("/:id", s.DeleteById)
	}
}

// FindAll fetches all strings or all strings associated with a certain
// thread if you pass a uuid as a `thread` query param
func (s *StringController) FindAll(c *gin.Context) {
	thread := c.Query("thread")
	if thread == "" {
		strings, err := s.Interactor.FindAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, strings)
		return
	}
	threadId, err := uuid.FromString(thread)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	strings, err := s.Interactor.FindAllByThread(threadId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, strings)
}

func (s *StringController) CreateOne(c *gin.Context) {
	var coreString core.String
	if err := c.ShouldBindJSON(&coreString); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Failed to bind: %s", err.Error()))
		return
	}
	newString, err := s.Interactor.CreateOne(coreString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, newString)
}

func (s *StringController) DeleteById(c *gin.Context) {
	stringId, err := uuid.FromString(c.Param("id"))
	if err != nil {
		log.Fatalf("failed to parse UUID %q: %v", s, err)
	}

	err = s.Interactor.DeleteById(stringId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, true)
}
