package main

import (
	"fmt"
	"html/template"

	"main/routes"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/joho/godotenv"
)

var store cookie.Store
var router *gin.Engine

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}

	store = cookie.NewStore([]byte(os.Getenv("SECRET_KEY")))
	store.Options(sessions.Options{MaxAge: 60 * 60 * 24, HttpOnly: true}) // expire in 24 hours and disable cookie access from js

}

func main() {

	// getting port env variable for render
	var host = os.Getenv("PORT")
	if host == "" {
		print("Website running on port 3000, go to localhost:3000\n")
		host = "3000"
	}
	// setting router & template for requests and static folders/files
	t := template.Must(template.ParseGlob("templates/*.html"))
	template.Must(t.ParseGlob("templates/partials/*.html"))
	router = gin.Default()
	router.SetHTMLTemplate(t)

	// setting static, sessions and auth middleware
	router.Static("/public/", "./public/")
	router.SetTrustedProxies(nil)
	router.Use(sessions.Sessions("mysession", store))
	router.Use(AuthRequired())

	// setting up routes to their controllers
	userMainRouterGroup := router.Group("/")
	postRouterGroup := router.Group("/posts")

	routes.SetupUserRoutes(userMainRouterGroup)
	routes.SetupPostRoutes(postRouterGroup)

	router.NoRoute(func(ctx *gin.Context) {
		val := sessions.Default(ctx).Get("user")
		ctx.HTML(http.StatusNotFound, "error.html", gin.H{
			"code":    http.StatusNotFound,
			"message": ctx.Request.URL.String() + " could not be found.",
			"User":    val,
		})
	})

	router.Run(":" + host)
}

// basic Page struct for info and user
type Page struct {
	Title  string
	Body   string
	Sample string
	User   interface{}
}

// middleware to check user authorized to perform action
func AuthRequired() gin.HandlerFunc {
	// returns the context handler function
	return func(ctx *gin.Context) {
		if ctx.Request.URL.String() != "/logout" && ctx.Request.URL.String() != "/posts/new" {
			ctx.Next()
			return
		}
		fmt.Println("Auth middleware on.")
		sess := sessions.Default(ctx)
		val := sess.Get("user")
		if val == nil {
			fmt.Println("not logged in to perfom action")
			ctx.HTML(http.StatusUnauthorized, "error.html", gin.H{
				"code":    401,
				"message": "Not authorized to perform action",
			})
			ctx.Abort()
			return
		}
		fmt.Println(val)
		fmt.Println("User authorized, middleware off.")
		ctx.Next()
	}
}
