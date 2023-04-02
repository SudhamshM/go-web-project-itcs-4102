package routes

import (
	"main/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	postController := controllers.PostController{}
	postRoutes := router.Group("/posts")
	{
		postRoutes.GET("", postController.GetPost)
		// postRoutes.POST("", postController.CreatePost)
		// postRoutes.PUT("/:id", postController.UpdatePost)
		// postRoutes.DELETE("/:id", postController.DeletePost)
	}
}
