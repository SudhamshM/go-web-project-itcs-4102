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
	fmt.Println("Getting post.")
	setupPostDB()
	objectID, idErr := primitive.ObjectIDFromHex(c.Param("id"))
	if idErr != nil {
		panic(idErr)
	}

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
		fmt.Println("Post found.")
		//condtional statement checking if the post userID is equal to the session
		//if post.userID == sessionuserID.. variable =true
		var userCondtion bool = false
		if post.UserID == sessions.Default(c).Get("user") {
			userCondtion = true

		}

		c.HTML(http.StatusOK, "post.html", gin.H{
			"post": post,
			"id":   post.ID.Hex(),
			"User": userCondtion,
		})
	}
}

func (u *PostController) CreatePost(ctx *gin.Context) {

	var r = ctx.Request
	setupPostDB()
	val := sessions.Default(ctx).Get("user")

	newPost := models.Post{
		Name:    r.FormValue("firstName"),
		Title:   r.FormValue("blogTitle"),
		Content: r.FormValue("blogContent"),
		ID:      primitive.NewObjectID(),
		UserID:  val,
	}
	_, insErr := postsCollection.InsertOne(ctx, newPost)
	if insErr != nil {
		panic(insErr)
	}

	cur, findErr := postsCollection.Find(ctx, bson.M{})
	if findErr != nil {
		panic(findErr)
	}
	var result []models.Post
	for cur.Next(ctx) {
		var post models.Post
		cur.Decode(&post)
		result = append(result, post)
	}
	ctx.HTML(http.StatusOK, "posts.html", gin.H{
		"error":    false,
		"bigArray": result,
		"hasPosts": true,
		"User":     val,
	})
}

func (u *PostController) NewPost(ctx *gin.Context) {
	val := sessions.Default(ctx).Get("user")
	ctx.HTML(http.StatusOK, "newblog.html", gin.H{
		"Title": "Create A Post",
		"User":  val,
	})
}

func (u *PostController) EditPost(ctx *gin.Context) {
	id := ctx.Param("id")
	var post models.Post
	objectID, _ := primitive.ObjectIDFromHex(id)
	setupPostDB()
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
	val := sessions.Default(ctx).Get("user")
	ctx.HTML(http.StatusOK, "edit.html", gin.H{
		"post": post,
		"User": val,
	})
}
func (u *PostController) UpdatePost(ctx *gin.Context) {
	title := ctx.PostForm("title")
	body := ctx.PostForm("body")
	id := ctx.Param("id")
	setupPostDB()
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

	update := bson.M{"$set": bson.M{"title": title, "content": body}}
	postsCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)

	// redirect them to the post they just edited
	ctx.Redirect(302, "/posts/"+id)
}
func (u *PostController) ViewPosts(ctx *gin.Context) {
	setupPostDB()
	cur, findErr := postsCollection.Find(ctx, bson.M{})
	if findErr != nil {
		panic(findErr)
	}
	var result []models.Post
	for cur.Next(ctx) {
		var post models.Post
		cur.Decode(&post)
		result = append(result, post)
	}
	val := sessions.Default(ctx).Get("user")
	ctx.HTML(http.StatusOK, "posts.html", gin.H{
		"error":    false,
		"bigArray": result,
		"hasPosts": true,
		"User":     val,
	})
}
func (u *PostController) DeletePost(ctx *gin.Context) {
	id := ctx.Param("id")

	var post models.Post
	objectID, _ := primitive.ObjectIDFromHex(id)
	setupPostDB()
	filter := bson.M{"_id": objectID}
	postsCollection.FindOne(ctx, filter).Decode(&post)
	if post.ID == primitive.NilObjectID {
		// if post is not there
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if post.UserID != sessions.Default(ctx).Get("user") {

		ctx.JSON(http.StatusNotFound, gin.H{"error": "User authorization failed"})
		return
	}

	postsCollection.DeleteOne(ctx, filter)

	ctx.Redirect(302, "/posts")

}

func setupPostDB() {
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
}
