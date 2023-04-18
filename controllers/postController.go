package controllers

import (
	"context"
	"fmt"
	"main/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
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
		//condtional statement checking if the post userID is equal to the session
		//if post.userID == sessionuserID.. variable =ture
		var userCondtion bool = false
		if post.UserID == sessions.Default(c).Get("user") {
			userCondtion = true

		}

		c.HTML(http.StatusOK, "post.html", gin.H{
			"post": post,
			"id":   post.ID.Hex(),
			"user": userCondtion,
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

func (u *PostController) EditPost(ctx *gin.Context) {
	id := ctx.Param("id")
	var post models.Post
	objectID, _ := primitive.ObjectIDFromHex(id)
	postsCollection := client.Database("goDatabase").Collection("posts")
	postsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post)
	if post.ID == primitive.NilObjectID {
		// if post is not there
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	// var userCondition bool = false
	if post.UserID != sessions.Default(ctx).Get("user") {
			
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User authorization failed"})
		return
	}

	ctx.HTML(http.StatusOK, "edit.html", gin.H{
		"post": post,
	})
}
func (u *PostController) UpdatePost(ctx *gin.Context) {
	title := ctx.PostForm("title")
	body := ctx.PostForm("body")
	id := ctx.Param("id")

	postsCollection := client.Database("goDatabase").Collection("posts")
	objectID, _ := primitive.ObjectIDFromHex(id)
	var post models.Post
	postsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post)
	if post.ID == primitive.NilObjectID {
		// if post is not there
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	if post.UserID != sessions.Default(ctx).Get("user") {
			
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User authorization failed"})
		return
	}


	update := bson.M{"$set": bson.M{"title": title,"content": body}}
	postsCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
		
	// redirect them to the post they just edited
	ctx.Redirect(302, "/posts/" + id)
}

