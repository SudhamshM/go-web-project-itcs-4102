package routes

import (
	"main/controllers"

	"github.com/gin-gonic/gin"
)

func SetupPostRoutes(rg *gin.RouterGroup) {
	postController := controllers.PostController{}
	rg.GET("/:id", postController.GetPost)
	rg.POST("/", postController.CreatePost)
	// postRoutes.DELETE("/:id", postController.DeletePost)
}

func SetupUserRoutes(rg *gin.RouterGroup) {
	userController := controllers.UserController{}
	rg.GET("/logout", userController.CreateUser)
}
