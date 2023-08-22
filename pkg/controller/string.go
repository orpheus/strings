package controller

import (
	"errors"
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
	UpdateName(stringId uuid.UUID, name string) error
	UpdateOrder(stringOrders []core.StringOrder) error
	DeleteById(id uuid.UUID) error
}

// RegisterRoutes creates a gin route grouping for the `/string` routes
func (s *StringController) RegisterRoutes(router *gin.RouterGroup) {
	skill := router.Group("/string")
	{
		skill.GET("", s.FindAll)
		skill.POST("", s.CreateOne)
		skill.DELETE("/:id", s.DeleteById)
		skill.PUT("/updateName", s.UpdateName)
		skill.PUT("/updateOrder", s.UpdateOrder)
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

// CreateOne creates a single string object
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

// UpdateName takes a string `id` and `name` as query params to update
// the name of that string.
func (s *StringController) UpdateName(c *gin.Context) {
	stringId, err := uuid.FromString(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Failed to parse string id: %s", err.Error()))
		return
	}

	newStringName := c.Query("name")

	err = s.Interactor.UpdateName(stringId, newStringName)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, true)
}

// StringOrderDTO is needed to bind the request body. The core.StringOrder struct defines
// `Id` as an uuid.UUID which gin cannot bind. So I needed to create a DTO struct to bind
// to and then convert to the core struct
type StringOrderDTO struct {
	Id    string `json:"id" binding:"required"`
	Order int    `json:"order"`
}

// UpdateOrder takes a list of id and order values and updates the strings
func (s *StringController) UpdateOrder(c *gin.Context) {
	var stringOrderDTOs []StringOrderDTO
	if err := c.ShouldBind(&stringOrderDTOs); err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New(fmt.Sprintf("Failed to bind `stringOrders`: %s", err.Error())))
		return
	}

	var stringOrders []core.StringOrder
	for _, s := range stringOrderDTOs {
		id, _ := uuid.FromString(s.Id)
		stringOrders = append(stringOrders, core.StringOrder{
			Id:    id,
			Order: s.Order,
		})
	}

	err := s.Interactor.UpdateOrder(stringOrders)

	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to update string order: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, true)
}

// DeleteById deletes a string by a string uuid passed as part
// of the req uri path
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
