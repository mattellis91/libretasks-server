package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
	"github.com/mattellis91/libretasks-server/models"
)

var Db *mongo.Database

func InitDb() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URL' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
  	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	Db := client.Database("dev")
	
	coll := Db.Collection("user")
	doc := models.User{Username: "test user", Email: "test@email.com"}

	res, err := coll.InsertOne(context.TODO(), doc)

	if err != nil {
		log.Fatal(err)
		log.Fatal("Failed to insert")
	}

	fmt.Printf("Inserted document with _id: %v\n", res.InsertedID)


	fmt.Println(Db.Name())
	fmt.Println(coll.Name())

	count, _ := coll.CountDocuments(context.TODO(), bson.M{})

	fmt.Println(count)

}