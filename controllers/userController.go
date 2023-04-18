package controllers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
}

func (u *UserController) SignupUser(ctx *gin.Context) {
	// Logic for creating a new user

	name := ctx.PostForm("username")
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")

	var currentSess sessions.Session = sessions.Default(ctx)

	result := Users{}
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
	user1 := Users{
		ID: primitive.NewObjectID(), Username: name, Email: email, Password: string(hash[:]),
	}

	_, err3 := usersCollection.InsertOne(ctx, user1)
	if err3 != nil {
		fmt.Println(err3)
		return
	}
	// showing success flash message

	currentSess.AddFlash("Account successfully created", "success")
	ok := currentSess.Flashes("success")
	currentSess.Flashes()
	currentSess.Save()
	ctx.HTML(http.StatusOK, "main.html", gin.H{
		"Title":       "Hello there",
		"Name":        name,
		"Body":        "Welcome to the UNC Charlotte Blog Website.",
		"Sample":      "Students can ask their peers for any help or share any advice for their peers relating to matters such as classes, clubs, sports, or other extracurricular activities.",
		"successMsgs": ok,
	})
}

func (u *UserController) StartLogin(ctx *gin.Context) {
	// Logic for creating a new user

	data := Page{
		Title: "Login",
		Body:  "Welcome to the login page",
	}
	ctx.HTML(http.StatusOK, "login.html", data)
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
			sess, _ := store.Get(ctx.Request, "mysession")
			sess.Values["user"] = user.ID.String()
			sess.AddFlash("You have successfully logged in!", "success")
			sess.Save(ctx.Request, ctx.Writer)
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
