package ginserver

import (
	"github.com/gin-gonic/gin"
	"github.com/orpheus/strings/pkg/controller/stringcon"
	"github.com/orpheus/strings/pkg/controller/threadcon"
	"github.com/orpheus/strings/pkg/infra/sqldb"
)

func Construct(r *gin.Engine, store *sqldb.Store) {
	v1Router := r.Group("/v1")
	v1Router.GET("/health", func(c *gin.Context) {
		c.JSON(200, "v1 healthy")
	})

	threadcon.NewThreadController(v1Router, store)
	stringcon.NewStringController(v1Router, store)
}
