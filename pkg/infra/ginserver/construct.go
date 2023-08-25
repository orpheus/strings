package ginserver

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/orpheus/strings/pkg/controller"
)

func Construct(r *gin.Engine, conn *pgxpool.Pool) {
	v1Router := r.Group("/v1")
	v1Router.GET("/v1/health", func(c *gin.Context) {
		c.JSON(200, "v1 healthy")
	})

	controller.NewThreadController(v1Router)
}
