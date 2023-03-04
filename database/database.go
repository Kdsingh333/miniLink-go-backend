package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
)

var DB *mongo.Collection
var ctx = context.TODO();
func Setup()*mongo.Collection {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")

	}
	Mongo_Atlas := os.Getenv("Mongo_Atlas")

	DB_NAME  := os.Getenv("DB_NAME")
	DB_COLLECTION_NAME := os.Getenv("DB_COLLECTION_NAME");
	
    clientOptions:= options.Client().ApplyURI(Mongo_Atlas);
	
	log.Println(1);

	client,err := mongo.Connect(ctx,clientOptions);
	log.Print(2);
	if err != nil{
		log.Fatal(err);
	}

	err = client.Ping(ctx,nil)
	log.Print(3);
	if err != nil{
		log.Fatal(err);
	}
    log.Print(4);
	DB= client.Database(DB_NAME).Collection(DB_COLLECTION_NAME );
    
    log.Print("DB connected");
	return DB;
}

// func GetDB() *mongo.Collection{
	
// 	return DB;
// }