package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var CustomerCollection *mongo.Collection
var UserCollection *mongo.Collection

func init() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
	uri := os.Getenv("uri")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	defer cancel()
	if err != nil {
		panic(err)
	}
	fmt.Println("Database Connected Successfully")

	CustomerCollection = client.Database(os.Getenv("dbname")).Collection(os.Getenv("customer_col"))
	UserCollection = client.Database(os.Getenv("dbname")).Collection(os.Getenv("user_col"))
	fmt.Println("Database is ready for operation.")
}
