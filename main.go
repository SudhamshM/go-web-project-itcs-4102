package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func main() {
	var host = os.Getenv("PORT")
	if host == "" {
		print("Website running on port 3000, go to localhost:3000\n")
		host = "3000"
	}
	blogArray = []BlogPosts{
		{FirstName: "Brijesh", TitlePost: "test title", ContentPost: "This is a test of posting a post."},
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
	router.Static("/public/", "./")

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

var blogArray []BlogPosts
var bigArray []BlogPosts

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Main handler: %v\n", r.URL.Path)
	data := Page{
		Title:  "Hello there",
		Body:   "Welcome to the UNC Charlotte Blog Website.",
		Sample: "Students can ask their peers for any help or share any advice for their peers relating to matters such as classes, clubs, sports, or other extracurricular activities.",
	}

	t, err := template_getter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err2 := t.ExecuteTemplate(w, "main.html", data)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

}

func about(w http.ResponseWriter, r *http.Request) {
	data := Page{
		Title:  "About Page!",
		Body:   "Welcome to my about page.",
		Sample: "ABOUT!",
	}
	t, err := template_getter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err2 := t.ExecuteTemplate(w, "about.html", data)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

}

func contact(w http.ResponseWriter, r *http.Request) {
	data := Page{
		Title:  "Contact Page",
		Body:   "Welcome to the contact page",
		Sample: "Please don't contact us about this site no one will response. ",
	}

	t, err := template_getter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err2 := t.ExecuteTemplate(w, "contact.html", data)
	if err2 != nil {
		fmt.Println(err2)
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

}

func signup(w http.ResponseWriter, r *http.Request) {
	data := Page{
		Title: "Sign Up",
		Body:  "Welcome to the sign up page",
	}

	t, err := template_getter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err2 := t.ExecuteTemplate(w, "signup.html", data)
	if err2 != nil {
		fmt.Println(err2)
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

}

func edit(w http.ResponseWriter, r *http.Request) {
	data := Page{
		Title:  "Edit Page!",
		Body:   "Welcome to my about page.",
		Sample: "Edit",
	}
	t, err := template_getter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err2 := t.ExecuteTemplate(w, "edit.html", data)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

}

func template_getter() (*template.Template, error) {
	t := template.New("")
	// Get a list of all files that match the "templates/*" pattern
	files, err := filepath.Glob("templates/*")
	if err != nil {
		print(err)
	}

	// Parse each file using the ParseFiles method of the template set
	t, err = t.ParseFiles(files...)
	if err != nil {
		print(err)
	}

	// makes sure template files are processed correctly
	return t, nil
}

func blog(w http.ResponseWriter, r *http.Request) {

	t, err := template_getter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err2 := t.ExecuteTemplate(w, "posts.html", bigArray)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

}

func newblog(w http.ResponseWriter, r *http.Request) {
	t, err := template_getter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// check if request method is GET or POST and return aptly
	if r.Method == "GET" {
		err2 := t.ExecuteTemplate(w, "newblog.html", nil)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == "POST" {
		var newBlog BlogPosts = BlogPosts{
			FirstName:   r.FormValue("firstName"),
			TitlePost:   r.FormValue("blogTitle"),
			ContentPost: r.FormValue("blogContent"),
			PostID:      uuid.New(),
		}
		bigArray = append(bigArray, newBlog)
		err2 := t.ExecuteTemplate(w, "posts.html", bigArray)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
			return
		}
	}

}

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
