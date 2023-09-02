package stringcon

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/orpheus/strings/pkg/controller/errorhandler"
	"github.com/orpheus/strings/pkg/infra/sqldb"
	"github.com/orpheus/strings/pkg/persistence/dao/stringdao"
	"github.com/orpheus/strings/pkg/persistence/repo/pgrepo/stringrepo"
	"github.com/orpheus/strings/pkg/service/stringsvc"
	"net/http"
)

func NewStringController(router *gin.RouterGroup, store *sqldb.Store) *StringController {
	stringDao := &stringdao.StringDao{Store: store}
	versionedStringDao := &stringdao.VersionedStringDao{Store: store}

	controller := &StringController{
		StringsService: stringsvc.NewStringService(
			stringrepo.NewStringRepository(stringDao, versionedStringDao),
		),
	}

	controller.RegisterRoutes(router)

	return controller
}

type StringController struct {
	StringsService StringsService
}

type StringsService interface {
	ArchiveString(stringId uuid.UUID) error
	RestoreString(stringId uuid.UUID) error
	ActivateString(stringId uuid.UUID) error
	DeactivateString(stringId uuid.UUID) error
	DeleteString(stringId uuid.UUID) error
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
	stringId := c.Param("id")
	stringIdUuid, err := uuid.Parse(stringId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorhandler.NewApiError(0, fmt.Sprintf("%s", err)))
		return
	}

	err = s.StringsService.DeleteString(stringIdUuid)
	if err != nil {
		errorhandler.HandleApiError(c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
