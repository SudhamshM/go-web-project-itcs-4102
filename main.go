package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"main/controllers"

	"main/routes"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/joho/godotenv"
)

var client *mongo.Client
var usersCollection *mongo.Collection
var store cookie.Store
var router *gin.Engine

// controllers
var postCtrl controllers.PostController
var userCtrl controllers.UserController

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}

	clientOptions := options.Client().
		ApplyURI(os.Getenv("DB_URL"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err2 error
	client, err2 = mongo.Connect(ctx, clientOptions)
	if err2 != nil {
		log.Fatal(err)
	}
	usersCollection = client.Database("goDatabase").Collection("users")

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
	router.Static("/public/", "./public/")
	router.SetTrustedProxies(nil)
	router.Use(sessions.Sessions("mysession", store))

	postCtrl = controllers.PostController{}
	userCtrl = controllers.UserController{}

	postRoutes := router.Group("/john")
	userRoutes := router.Group("/user")
	routes.SetupPostRoutes(postRoutes)
	routes.SetupUserRoutes(userRoutes)

	router.Use(func(ctx *gin.Context) {
		if ctx.Request.URL.String() != "/logout" && ctx.Request.URL.String() != "/posts/new" {
			ctx.Next()
			return
		}
		fmt.Println("auth middleware on")
		sess, _ := store.Get(ctx.Request, "mysession")
		val, ok := sess.Values["user"]
		if !ok {
			fmt.Println("not logged in to perfom action")
			ctx.HTML(http.StatusUnauthorized, "error.html", gin.H{
				"code":    401,
				"message": "Not authorized to perform action",
			})
			ctx.Abort()
			return
		}
		fmt.Println(val)
		fmt.Println("user authorized")
		fmt.Println("middleware off")
		ctx.Next()
	})

	// route handlers
	router.GET("/", func(ctx *gin.Context) {
		success := sessions.Default(ctx).Flashes("success")
		errMsgs := sessions.Default(ctx).Flashes("error")

		// clearing the flash before rendering
		sessions.Default(ctx).Flashes()
		sessions.Default(ctx).Save()
		// either object id string or nil
		val := sessions.Default(ctx).Get("user")

		ctx.HTML(http.StatusOK, "main.html", gin.H{
			"Title":       "Hello there",
			"Body":        "Welcome to the UNC Charlotte Blog Website.",
			"Sample":      "Students can ask their peers for any help or share any advice for their peers relating to matters such as classes, clubs, sports, or other extracurricular activities.",
			"successMsgs": success,
			"errorMsgs":   errMsgs,
			"user":        val,
		})
	})

	// post routes

	router.GET("/posts", postCtrl.ViewPosts)
	router.GET("/posts/:id", postCtrl.GetPost)
	router.GET("/posts/new", postCtrl.NewPost)
	router.POST("/posts", postCtrl.CreatePost)
	router.GET("/edit/:id", postCtrl.EditPost)
	router.POST("/edit/:id", postCtrl.UpdatePost)
	router.POST("/delete/:id", postCtrl.DeletePost)


	// user routes

	router.GET("/login", userCtrl.StartSignup)
	router.POST("/login", userCtrl.LoginUser)
	router.GET("/signup", userCtrl.StartSignup)
	router.POST("/signup", userCtrl.SignupUser)
	router.GET("/logout", userCtrl.LogoutUser)

	router.GET("/about", func(ctx *gin.Context) {
		data := Page{
			Title:  "About Page!",
			Body:   "Welcome to my about page.",
			Sample: "ABOUT!",
		}
		ctx.HTML(http.StatusOK, "about.html", data)
	})

	router.GET("/contact", func(ctx *gin.Context) {
		data := Page{
			Title:  "Contact Page",
			Body:   "Welcome to the contact page",
			Sample: "Please don't contact us about this site no one will response. ",
		}
		ctx.HTML(http.StatusOK, "contact.html", data)
	})





	router.NoRoute(func(ctx *gin.Context) {

		ctx.HTML(http.StatusNotFound, "error.html", gin.H{
			"code":    http.StatusNotFound,
			"message": ctx.Request.URL.String() + " could not be found.",
		})
	})

	router.Run(":" + host)
}

// basic Page struct for info
type Page struct {
	Title  string
	Body   string
	Sample string
}

// blog post struct
type BlogPosts struct {
	FirstName   string `json:"firstname"`
	TitlePost   string `json:"title"`
	ContentPost string `json:"contentpost"`
	PostID      uuid.UUID
}

type Users struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

func getUserByEmail(ctx *gin.Context, email string) *Users {
	var result = Users{}
	usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&result)
	if result.ID == primitive.NilObjectID {
		return nil
	}
	return &result
}
