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
	"golang.org/x/crypto/bcrypt"

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

	router.GET("/login", func(ctx *gin.Context) {
		data := Page{
			Title: "Login",
			Body:  "Welcome to the login page",
		}
		ctx.HTML(http.StatusOK, "login.html", data)
	})

	router.POST("/login", func(ctx *gin.Context) {
		email := ctx.PostForm("email")
		password := ctx.PostForm("password")

		user := getUserByEmail(ctx, email)

		if user == nil {
			fmt.Println("user not found with email")
			ctx.HTML(http.StatusOK, "error.html", gin.H{
				"code":    404,
				"message": "User not found with given email",
			})
			return
		} else {
			fmt.Println("user found")
			fmt.Println(user)
			pwdCheck := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if pwdCheck == nil {
				fmt.Println("login success")
				// store user id in session
				sess, _ := store.Get(ctx.Request, "mysession")
				sess.Values["user"] = user.ID.String()
				sess.AddFlash("You have successfully logged in!", "success")
				sess.Save(ctx.Request, ctx.Writer)
				ctx.Redirect(302, "/")
				return
			} else {
				fmt.Println("wrong password", pwdCheck)
				ctx.HTML(http.StatusOK, "error.html", gin.H{
					"code":    401,
					"message": "Incorrect password",
				})
				return
			}
		}
	})

	router.GET("/signup", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "signup.html", gin.H{
			"Title":        "Sign Up",
			"Body":         "Welcome to the sign up page",
			"error":        nil,
			"errorMessage": nil,
		})
	})

	router.GET("/logout", func(ctx *gin.Context) {
		// sess, _ := store.Get(ctx.Request, "mysession")
		// sess.Values["user"] = nil

		// to use above logic, update auth middleware to check for nil instead of ok
		sessions.Default(ctx).Clear()
		sessions.Default(ctx).AddFlash("You have successfully logged out!", "success")
		sessions.Default(ctx).Save()
		fmt.Println("logged out")
		ctx.Redirect(302, "/")
	})

	router.POST("/signup", func(ctx *gin.Context) {
		name := ctx.PostForm("username")
		email := ctx.PostForm("email")
		password := ctx.PostForm("password")

		var currentSess sessions.Session = sessions.Default(ctx)

		result := Users{}
		usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&result)
		// check if user doesn't exist in db already
		if result.Email == email {
			ctx.HTML(http.StatusBadRequest, "signup.html", gin.H{
				"Title":        "Sign Up",
				"Body":         "Welcome to the sign up page",
				"error":        true,
				"errorMessage": "Email already in use.",
			})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		if err != nil {
			fmt.Println(err)
			return
		}
		user1 := Users{
			ID: primitive.NewObjectID(), Username: name, Email: email, Password: string(hash[:]),
		}

		_, err3 := usersCollection.InsertOne(ctx, user1)
		if err3 != nil {
			fmt.Println(err3)
			return
		}
		// showing success flash message

		currentSess.AddFlash("Account successfully created", "success")
		ok := currentSess.Flashes("success")
		currentSess.Flashes()
		currentSess.Save()
		ctx.HTML(http.StatusOK, "main.html", gin.H{
			"Title":       "Hello there",
			"Name":        name,
			"Body":        "Welcome to the UNC Charlotte Blog Website.",
			"Sample":      "Students can ask their peers for any help or share any advice for their peers relating to matters such as classes, clubs, sports, or other extracurricular activities.",
			"successMsgs": ok,
		})
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
