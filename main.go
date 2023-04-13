package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"main/controllers"
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

	"github.com/joho/godotenv"
)

var client *mongo.Client

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
	user1 := Users{
		ID: primitive.NewObjectID(), Username: "Melvin Sudhamsh", Email: "b@a.com", Password: "good",
	}
	_, err3 := client.Database("goDatabase").Collection("users").InsertOne(ctx, user1)
	if err3 != nil {
		fmt.Println(err3)
	}
}

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

	postCtrl := controllers.PostController{}
	// switch to controller defined routes for future
	router.GET("/undefined", postCtrl.GetPost)
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

	router.POST("/signup", func(ctx *gin.Context) {
		name := ctx.PostForm("username")
		email := ctx.PostForm("email")
		password := ctx.PostForm("password")

		addedUser := databaseCollection.FindOne(context.Background(), bson.M{"username": name})

		if addedUser == nil {
			ctx.AbortWithStatus(500)
		}

		addedEmail := databaseCollection.FindOne(context.Background(), bson.M{"email": email})

		if addedEmail == nil {
			ctx.AbortWithStatus(500)
		}

		var user Users
		addedPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

		if addedPassword == nil {
			ctx.AbortWithStatus(500)
		}

		ctx.HTML(http.StatusOK, "main.html", gin.H{
			"Title":  "Hello there",
			"Name":   name,
			"Body":   "Welcome to the UNC Charlotte Blog Website.",
			"Sample": "Students can ask their peers for any help or share any advice for their peers relating to matters such as classes, clubs, sports, or other extracurricular activities.",
		})

	})

	// (Aiden) editing a page:

	router.GET("/edit/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		post := getPostById(id)
		fmt.Println("finding post...")
		if post == nil {
			// if post is not there
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}

		ctx.HTML(http.StatusOK, "edit.html", gin.H{
			"post": post,
		})
	})

	router.POST("/edit/:id", func(ctx *gin.Context) {
		title := ctx.PostForm("title")
		body := ctx.PostForm("body")
		id := ctx.Param("id")
		post := getPostById(id)

		fmt.Println("finding post...")
		if post == nil {
			// if post is not there
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		post.TitlePost = title
		post.ContentPost = body

		// redirect them to the post they just edited
		ctx.HTML(http.StatusOK, "post.html", gin.H{
			"post": post,
		})
	})

	//incomplete delete post method
	//need to loop through the posts and derive the needed index
	//delete the index from the bigarray
	router.POST("/delete/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var nonsense []BlogPosts

		//loop through each value in bigArray
		for _, v := range bigArray {
			//compare the postID to the given postID
			if v.PostID.String() == id {

			} else {
				nonsense = append(nonsense, v)
			}
		}
		bigArray = nonsense

		// redirect them to the post they just edited
		ctx.Redirect(302, "/posts")

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
