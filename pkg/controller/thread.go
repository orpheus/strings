package controller

import (
	"github.com/gin-gonic/gin"
)

type ThreadController struct {
	ThreadService ThreadService
}

type ThreadService interface {
}

func (t *ThreadController) RegisterRoutes(router *gin.RouterGroup) {
	threads := router.Group("/threads")
	{
		threads.POST("", t.PostThreads)
		threads.GET("", t.PostThreads)
	}

	thread := router.Group("/thread")
	{
		thread.GET("", t.Archive)
		thread.POST("/archive/:id", t.Archive)
		thread.POST("/restore/:id", t.Restore)
		thread.POST("/activate/:id", t.Activate)
		thread.POST("/deactivate/:id", t.Deactivate)
		thread.POST("/delete/:id", t.Delete)
	}
}

func (t *ThreadController) PostThreads(c *gin.Context) {

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
