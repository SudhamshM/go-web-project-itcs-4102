package controllers

import (
	"context"
	"fmt"
	"main/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var postsCollection *mongo.Collection

type PostController struct {
}

func (u *PostController) GetPost(c *gin.Context) {
	// Logic for creating a new user
	fmt.Println("getting post")
	objectID, idErr := primitive.ObjectIDFromHex(c.Param("id"))
	if idErr != nil {
		panic(idErr)
	}

	var DB_URL string = os.Getenv("DB_URL")
	clientOptions := options.Client().
		ApplyURI(DB_URL)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err2 error
	client, err2 = mongo.Connect(ctx, clientOptions)
	if err2 != nil {
		panic(err2)
	}
	postsCollection = client.Database("goDatabase").Collection("posts")
	var post models.Post
	postsCollection.FindOne(c, bson.M{"_id": objectID}).Decode(&post)
	fmt.Println(post.Title, post.ID)

	if post.ID == primitive.NilObjectID {
		fmt.Println("No post found.")
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"code":    http.StatusNotFound,
			"message": "Post not found",
		})
		c.Abort()
		return
	} else {
		fmt.Println("Post found")

		c.HTML(http.StatusOK, "post.html", gin.H{
			"post": post,
			"id":   post.ID.Hex(),
		})
	}
}

func (u *PostController) CreatePost(c *gin.Context) {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	var DB_URL string = os.Getenv("DB_URL")
	clientOptions := options.Client().
		ApplyURI(DB_URL)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err2 error
	client, err2 = mongo.Connect(ctx, clientOptions)
	if err2 != nil {
		panic(err2)
	}
	postsCollection = client.Database("goDatabase").Collection("posts")
	newPost := models.Post{
		Name:    "John Wick",
		Title:   "",
		Content: "Go watch movie!",
		ID:      primitive.NewObjectID(),
	}
	_, insErr := postsCollection.InsertOne(c, newPost)
	if insErr != nil {
		panic(insErr)
	}
}
