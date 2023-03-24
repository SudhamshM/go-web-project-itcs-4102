package main

import (
	"fmt"
	"net/http"
	"os"

	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	var databaseCollection *mongo.Collection
	ctx := context.TODO()
	options := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, options)

	if err != nil {
		panic(err)
	}

	databaseCollection = client.Database("goDatabase").Collection("users")
	
	// getting port env variable for render
	var host = os.Getenv("PORT")
	if host == "" {
		print("Website running on port 3000, go to localhost:3000\n")
		host = "3000"
	}
	// setting router & template for requests and static folders/files
	t := template.Must(template.ParseGlob("templates/*.html"))
	template.Must(t.ParseGlob("templates/partials/*.html"))
	router := gin.Default()
	router.SetHTMLTemplate(t)
	router.Static("/public/", "./public/")
	router.SetTrustedProxies(nil)

	// route handlers
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "main.html", gin.H{
			"Title":  "Hello there",
			"Body":   "Welcome to the UNC Charlotte Blog Website.",
			"Sample": "Students can ask their peers for any help or share any advice for their peers relating to matters such as classes, clubs, sports, or other extracurricular activities.",
		})
	})

	router.GET("/posts", func(ctx *gin.Context) {
		if len(bigArray) == 0 {
			ctx.HTML(http.StatusOK, "posts.html", gin.H{
				"error":    true,
				"hasPosts": false,
			})
			return
		} else {
			ctx.HTML(http.StatusOK, "posts.html", gin.H{
				"error":    false,
				"bigArray": bigArray,
				"hasPosts": true,
			})
		}
	})

	router.GET("/posts/:id", singlePost)

	router.GET("/newblog", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "newblog.html", nil)
	})

	router.POST("/newblog", func(ctx *gin.Context) {
		var r = ctx.Request
		var newBlog BlogPosts = BlogPosts{
			FirstName:   r.FormValue("firstName"),
			TitlePost:   r.FormValue("blogTitle"),
			ContentPost: r.FormValue("blogContent"),
			PostID:      uuid.New(),
		}
		bigArray = append(bigArray, newBlog)
		ctx.HTML(http.StatusOK, "posts.html", gin.H{
			"error":    false,
			"bigArray": bigArray,
			"hasPosts": true,
		})
	})

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

	router.GET("/signup", func(ctx *gin.Context) {
		data := Page{
			Title: "Sign Up",
			Body:  "Welcome to the sign up page",
		}
		ctx.HTML(http.StatusOK, "signup.html", data)
	})
	
	router.POST("signup", func(ctx *gin.Context) {
		name := ctx.PostForm("username")
		email := ctx.PostForm("email")
		password := ctx.PostForm("password")

		addedUser := databaseCollection.FindOne(context.Background(), bson.M{"username": name})

		if addedUser != nil {
			ctx.AbortWithStatus(500)
		}

		addedEmail := databaseCollection.FindOne(context.Background(), bson.M{"email": email})

		if addedEmail != nil {
			ctx.AbortWithStatus(500)
		}

		var user Users
		addedPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

		if addedPassword != nil {
			ctx.AbortWithStatus(500)
		}

		//ctx.JSON(200, gin.H{"Access": fmt.Sprintf("Hello %s ", name)})
		ctx.HTML(http.StatusOK, "main.html", gin.H{
			"Title":  "Hello there",
			"Name":   name,
			"Body":   "Welcome to the UNC Charlotte Blog Website.",
			"Sample": "Students can ask their peers for any help or share any advice for their peers relating to matters such as classes, clubs, sports, or other extracurricular activities.",
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
	ID       primitive.ObjectID `bson:_id`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

// array to hold all posts
var bigArray []BlogPosts

func singlePost(c *gin.Context) {
	id := c.Param("id")
	post := getPostById(id)
	fmt.Println("finding post...")
	if post == nil {
		// if post is not there
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.HTML(http.StatusOK, "post.html", gin.H{
		"post": post,
	})

}

// get posts by id to show specific post
func getPostById(id string) *BlogPosts {
	for i := 0; i < len(bigArray); i++ {
		if bigArray[i].PostID.String() == id {
			return &bigArray[i]
		}
	}
	return nil
}
