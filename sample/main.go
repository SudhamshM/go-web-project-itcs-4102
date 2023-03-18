package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/about", about)
	http.ListenAndServe(":3000", nil)
}

type Page struct {
	Title  string
	Body   string
	Sample string
}

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

	return t, nil
}
