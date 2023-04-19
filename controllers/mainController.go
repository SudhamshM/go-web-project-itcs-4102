package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type MainController struct {
}

func (m *MainController) GetIndex(ctx *gin.Context) {
	val := sessions.Default(ctx).Get("user")
	success := sessions.Default(ctx).Flashes("success")
	errMsgs := sessions.Default(ctx).Flashes("error")

	// clearing the flash before rendering
	sessions.Default(ctx).Flashes()
	sessions.Default(ctx).Save()

	ctx.HTML(http.StatusOK, "main.html", gin.H{
		"Title":       "Hello there",
		"Body":        "Welcome to the UNC Charlotte Blog Website.",
		"Sample":      "Students can ask their peers for any help or share any advice for their peers relating to matters such as classes, clubs, sports, or other extracurricular activities.",
		"successMsgs": success,
		"errorMsgs":   errMsgs,
		"User":        val,
	})
}

func (m *MainController) GetAbout(ctx *gin.Context) {
	val := sessions.Default(ctx).Get("user")
	data := Page{
		Title:  "About Page!",
		Body:   "Welcome to my about page.",
		Sample: "ABOUT!",
		User:   val,
	}
	ctx.HTML(http.StatusOK, "about.html", data)
}

func (m *MainController) GetContact(ctx *gin.Context) {
	val := sessions.Default(ctx).Get("user")
	data := Page{
		Title:  "Contact Page",
		Body:   "Welcome to the contact page",
		Sample: "Please don't contact us about this site no one will response. ",
		User:   val,
	}
	ctx.HTML(http.StatusOK, "contact.html", data)
}
