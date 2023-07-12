package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID primitive.ObjectID `bson:"_id"`
	Name string `bson:"name"`
}

func InitDb() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URL' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	db := client.Database("libretasks")
	coll := db.Collection("user")

	fmt.Println(db.Name())
	fmt.Println(coll.Name())

	count, _ := coll.CountDocuments(context.TODO(), bson.M{})

	fmt.Println(count)

	var result User
	err = coll.FindOne(context.TODO(), bson.D{{"name", "test"}}).Decode(&result)   //coll.Find(context.TODO(), bson.D{{}}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		fmt.Printf("No documents were found")
		return
	}

	if err != nil {
		panic(err)
	}


	jsonData, err := json.MarshalIndent(result, "", "	")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", jsonData)

}