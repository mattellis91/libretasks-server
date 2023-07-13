package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mattellis91/libretasks-server/models"
)

var Db *mongo.Database = nil
var Connection *mongo.Client = nil

func InitDb() {
	
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_DB_NAME")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URL' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
  	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	c, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}

	Connection = c

	Db = c.Database(dbName)
	coll := Db.Collection("user")
	doc := models.User{ID: primitive.NewObjectID(), Username:  "test user 2", Email: "test2@email.com"}

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

func CloseConnection() {
	if Db != nil {
		Connection.Disconnect(context.TODO())
		Db = nil
	}
}