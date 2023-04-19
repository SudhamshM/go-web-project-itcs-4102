package routes

import (
	"main/controllers"

	"github.com/gin-gonic/gin"
)

func SetupPostRoutes(rg *gin.RouterGroup) {
	postController := controllers.PostController{}
	rg.GET("/:id", postController.GetPost)
	rg.POST("", postController.CreatePost)
	rg.GET("", postController.ViewPosts)
	rg.GET("/new", postController.NewPost)
	rg.GET("/:id/edit", postController.EditPost)
	rg.POST("/:id/edit", postController.UpdatePost)
	rg.POST("/:id/delete", postController.DeletePost)
}

func SetupUserRoutes(rg *gin.RouterGroup) {
	userController := controllers.UserController{}
	mainController := controllers.MainController{}

	// adding user routes
	rg.GET("/logout", userController.LogoutUser)
	rg.GET("/login", userController.StartLogin)
	rg.POST("/login", userController.LoginUser)
	rg.GET("/signup", userController.StartSignup)
	rg.POST("/signup", userController.SignupUser)

	// adding default main routes
	rg.GET("", mainController.GetIndex)
	rg.GET("/about", mainController.GetAbout)
	rg.GET("/contact", mainController.GetContact)
}
