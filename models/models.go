package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	Name    string             `bson:"name"`
	Title   string             `bson:"title"`
	Content string             `bson:"content"`
	ID      primitive.ObjectID `bson:"_id"`
}

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}
