package ginserver

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

func Construct(r *gin.Engine, conn *pgxpool.Pool) {
	//tmpLogger := &logging.TmpLogger{}
	//
	//v1Router := r.Group("/api")
	//v1Router.GET("/health", func(c *gin.Context) {
	//	c.JSON(200, "healthy")
	//})
	//
	//threadRepository := &thread.Repository{
	//	DB:     conn,
	//	Logger: tmpLogger,
	//}
	//
	//stringRepository := &string.StringRepository{
	//	DB:     conn,
	//	Logger: tmpLogger,
	//}
	//
	//threadController := &thread.Controller{
	//	Interactor: &system.ThreadInteractor{
	//		Repo:          threadRepository,
	//		StringDeleter: stringRepository,
	//		Logger:        tmpLogger,
	//	},
	//	Logger: tmpLogger,
	//}
	//
	//stringController := string.StringController{
	//	Interactor: &system.StringInteractor{
	//		StringRepository: stringRepository,
	//		Logger:           tmpLogger,
	//	},
	//	Logger: nil,
	//}
	//
	//threadController.RegisterRoutes(v1Router)
	//stringController.RegisterRoutes(v1Router)
}
