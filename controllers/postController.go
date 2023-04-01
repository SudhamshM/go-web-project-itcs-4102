package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type PostController struct {
}

func (u *PostController) GetPost(c *gin.Context) {
	// Logic for creating a new user
	id := c.Param("id")
	fmt.Println("finding post..." + id)
}
