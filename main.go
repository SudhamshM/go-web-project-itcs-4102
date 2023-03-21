package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func main() {
	var host = os.Getenv("PORT")
	if host == "" {
		print("Website running on port 3000, go to https://localhost:3000\n")
		host = "3000"
	}
	blogArray = []BlogPosts{
		{FirstName: "Brijesh", TitlePost: "test title", ContentPost: "This is a test of posting a post."},
	}

	router := gin.Default()
	router.GET("/newblog", func(c *gin.Context) {
		c.JSON(200, blogArray)
	})

	router.POST("/newblog", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, &blogArray)

	})

	// serving static files using file server
	fs := http.FileServer(http.Dir("public/"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	http.HandleFunc("/", handler)
	http.HandleFunc("/about", about)
	http.HandleFunc("/contact", contact)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/edit", edit)
	http.ListenAndServe(":"+host, nil)

	router.Run()
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
}

var blogArray []BlogPosts

func handler(w http.ResponseWriter, r *http.Request) {
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
	// fmt.Println(files)
	return t, nil
}

func blog(w http.ResponseWriter, r *http.Request) {
	/*blogArray = []BlogPosts{
		{FirstName: "Brijesh", TitlePost: "test title", ContentPost: "This is a test of posting a post."},
		//{FirstName: "Ajay", TitlePost: "Title 1", ContentPost: "This is another post."},
		//{FirstName: "Mevlin", TitlePost: "Title 2", ContentPost: "This is another post part two."},
	}*/
	blogArray := BlogPosts{
		FirstName:   "Brijesh",
		TitlePost:   "test title",
		ContentPost: "This is a test of posting a post.",
	}

	t, err := template_getter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err2 := t.ExecuteTemplate(w, "posts.html", blogArray)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

}

func newblog(w http.ResponseWriter, r *http.Request) {
	blogArrays := BlogPosts{
		FirstName:   "Brijesh",
		TitlePost:   "test title",
		ContentPost: "This is a test of posting a post.",
	}
	t, err := template_getter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err2 := t.ExecuteTemplate(w, "newblog.html", blogArrays)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

}
