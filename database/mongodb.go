package database

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)



func ConnectDB() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}


	url := os.Getenv("MONGO_URL"); 
	
	if url == "" {
		log.Fatal("MongoDB URL IS NOT SET")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))

	if err != nil {
		log.Fatal(err)
	}

	// check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected to mongodb")
	}

	return client
}

var Client *mongo.Client = ConnectDB()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("go-fiber-ecommerce").Collection(collectionName)
	return collection
}