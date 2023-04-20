package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"main/models"

	"github.com/gin-contrib/sessions"
)

type UserController struct {
}

var usersCollection *mongo.Collection

func (u *UserController) SignupUser(ctx *gin.Context) {
	// Logic for creating a new user

	name := ctx.PostForm("username")
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")

	currentSess := sessions.Default(ctx)
	// setup DB before proceeding
	setupUserDB()
	result := models.User{}
	usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&result)
	// check if user doesn't exist in db already
	if result.Email == email {
		ctx.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"Title":        "Sign Up",
			"Body":         "Welcome to the sign up page",
			"error":        true,
			"errorMessage": "Email already in use.",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		fmt.Println(err)
		return
	}
	user1 := models.User{
		ID: primitive.NewObjectID(), Username: name, Email: email, Password: string(hash[:]),
	}

	_, err3 := usersCollection.InsertOne(ctx, user1)
	if err3 != nil {
		fmt.Println(err3)
		return
	}
	// showing success flash message

	currentSess.AddFlash("Account successfully created", "success")
	currentSess.Flashes()
	currentSess.Save()
	ctx.Redirect(http.StatusFound, "/")
}

func (u *UserController) StartLogin(ctx *gin.Context) {
	// Logic for creating a new user
	val := sessions.Default(ctx).Get("user")
	// fmt.Println(val).... add check if user is already logged in to prevent from going here
	data := Page{
		Title: "Login",
		Body:  "Welcome to the login page",
		User:  val,
	}

	if val == nil {
		ctx.HTML(http.StatusOK, "login.html", data)

	} else {
		sessions.Default(ctx).AddFlash("You cannot login while already logging in", "warning")
		sessions.Default(ctx).Save()
		fmt.Println("You cannot login while already logging in")
		ctx.Redirect(302, "/")
	}

}

func (u *UserController) LogoutUser(ctx *gin.Context) {
	// Logic for creating a new user
	// sess, _ := store.Get(ctx.Request, "mysession")
	// sess.Values["user"] = nil

	// to use above logic, update auth middleware to check for nil instead of ok
	sessions.Default(ctx).Clear()
	sessions.Default(ctx).AddFlash("You have successfully logged out!", "success")
	sessions.Default(ctx).Save()
	fmt.Println("logged out")
	ctx.Redirect(302, "/")

}

func (u *UserController) LoginUser(ctx *gin.Context) {
	// Logic for creating a new user

	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	setupUserDB()

	user := getUserByEmail(ctx, email)

	if user == nil {
		fmt.Println("user not found with email")
		ctx.HTML(http.StatusOK, "error.html", gin.H{
			"code":    404,
			"message": "User not found with given email",
		})
		return
	} else {
		fmt.Println("user found")
		fmt.Println(user)
		pwdCheck := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if pwdCheck == nil {
			fmt.Println("login success")
			// store user id in session
			sess := sessions.Default(ctx)
			sess.Set("user", user.ID.String())
			sess.AddFlash("You have successfully logged in!", "success")
			sess.Save()
			ctx.Redirect(302, "/")
			return
		} else {
			fmt.Println("wrong password", pwdCheck)
			ctx.HTML(http.StatusOK, "error.html", gin.H{
				"code":    401,
				"message": "Incorrect password",
			})
			return
		}
	}
}

func (u *UserController) StartSignup(ctx *gin.Context) {
	// Logic for creating a new user

	ctx.HTML(http.StatusOK, "signup.html", gin.H{
		"Title":        "Sign Up",
		"Body":         "Welcome to the sign up page",
		"error":        nil,
		"errorMessage": nil,
	})

}

// basic Page struct for info
type Page struct {
	Title  string
	Body   string
	Sample string
	User   interface{}
}

func getUserByEmail(ctx *gin.Context, email string) *models.User {
	setupUserDB()
	var result = models.User{}
	usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&result)
	if result.ID == primitive.NilObjectID {
		return nil
	}
	return &result
}

func setupUserDB() {
	var DB_URL string = os.Getenv("DB_URL")
	clientOptions := options.Client().ApplyURI(DB_URL)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err2 error
	client, err2 = mongo.Connect(ctx, clientOptions)
	if err2 != nil {
		panic(err2)
	}

	usersCollection = client.Database("goDatabase").Collection("users")
}
