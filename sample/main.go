package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/about", about)
	http.ListenAndServe(":8080", nil)
}

type Page struct {
	Title  string
	Body   string
	Sample string
}

func handler(w http.ResponseWriter, r *http.Request) {
	data := Page{
		Title:  "Hello, World!",
		Body:   "Welcome to my Go web application.",
		Sample: "Ajay",
	}
	templates := template.Must(template.ParseFiles("navbar.html", "main.html"))
	// Execute the navbar template
	err := templates.ExecuteTemplate(w, "main.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func about(w http.ResponseWriter, r *http.Request) {
	data := Page{
		Title:  "About Page!",
		Body:   "Welcome to my about page.",
		Sample: "ABOUT!",
	}

	templates := template.Must(template.ParseFiles("about.html", "main.html", "navbar.html"))
	// Execute the navbar template
	err := templates.ExecuteTemplate(w, "about.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
