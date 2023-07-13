package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mattellis91/libretasks-server/db"
	"github.com/mattellis91/libretasks-server/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {

	dbc := db.GetDatabaseConnection()
	coll := dbc.Collection("user")
	doc := models.User{ID: primitive.NewObjectID(), Username:  "test user 2", Email: "test2@email.com"}

	res, err := coll.InsertOne(context.TODO(), doc)

	if err != nil {
		log.Fatal(err)
		log.Fatal("Failed to insert")
	}

	fmt.Printf("Inserted document with _id: %v\n", res.InsertedID)

	db.CloseConnection()
	
}
