package ginserver

import (
	"github.com/gin-gonic/gin"
	"github.com/orpheus/strings/pkg/controller"
	"github.com/orpheus/strings/pkg/infra/sqldb"
)

func Construct(r *gin.Engine, store *sqldb.Store) {
	v1Router := r.Group("/v1")
	v1Router.GET("/v1/health", func(c *gin.Context) {
		c.JSON(200, "v1 healthy")
	})

	controller.NewThreadController(v1Router, store)
}
