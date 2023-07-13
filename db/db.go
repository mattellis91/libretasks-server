package db

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Db *mongo.Database = nil
var Client *mongo.Client = nil

func createConnection() error {
	
	if err := godotenv.Load(); err != nil {
		panic("Could not load env")
	}
	
	dbName := os.Getenv("MONGODB_DB_NAME")
	
	//Client should always be nil if it gets here but just incase check if the client connection exists before making a new one
	if Client == nil {
		
		uri := os.Getenv("MONGODB_URI")
		
		if uri == "" {
			log.Fatal("You must set your 'MONGODB_URL' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
		}

		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

		c, err := mongo.Connect(context.TODO(), opts)

		if err != nil {
			panic(err)
		}

		Client = c

	}
	
	Db = Client.Database(dbName)

	return nil
}

func GetDatabaseConnection() *mongo.Database {
	if Db == nil {
		if error := createConnection(); error != nil {
			panic("Could not connect to database")
		}
	}
	return Db
} 

func CloseConnection() {
	if Client != nil {
		Client.Disconnect(context.TODO())
		Client = nil
		Db = nil
	}
}