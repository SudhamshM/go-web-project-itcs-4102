package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	var host = os.Getenv("PORT")
	if host == "" {
		print("Website running on port 3000, go to localhost:3000\n")
		host = "3000"
	}

	router := gin.Default()
	router.GET("/posts/:id", singlePost)

	// serving static files using file server
	// fs := http.FileServer(http.Dir("public/"))
	// http.Handle("/public/", http.StripPrefix("/public/", fs))
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", handler)
	// // http.HandleFunc("/", handler)
	// mux.HandleFunc("/about", about)
	// mux.HandleFunc("/contact", contact)
	// mux.HandleFunc("/signup", signup)
	// mux.HandleFunc("/edit", edit)
	// mux.HandleFunc("/blog", blog)
	// mux.HandleFunc("/newblog", newblog)
	// http.ListenAndServe(":"+host, mux)
	router.LoadHTMLGlob("templates/*")
	router.Static("/public/", "./public/")

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

	router.Run(":" + host)
}

type Page struct {
	Title  string
	Body   string
	Sample string
}

type BlogPosts struct {
	FirstName   string `json:"firstname"`
	TitlePost   string `json:"title"`
	ContentPost string `json:"contentpost"`
	PostID      uuid.UUID
}

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

func getPostById(id string) *BlogPosts {
	for i := 0; i < len(bigArray); i++ {
		if bigArray[i].PostID.String() == id {
			return &bigArray[i]
		}
	}
	return nil
}
