package code

import (
	"context"
	"fmt"
	"time"

	"github.com/EmeraldLS/phsps-api/db"
	"github.com/EmeraldLS/phsps-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMaxCustomerCode() int {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	filter := bson.M{}
	findOptions := options.Find().SetSort(bson.M{"customer_code": -1}).SetLimit(1)
	cursor, _ := db.CustomerCollection.Find(ctx, filter, findOptions)
	defer cursor.Close(ctx)
	var customers []model.Customer
	for cursor.Next(ctx) {
		var customer model.Customer
		cursor.Decode(&customer)
		customers = append(customers, customer)
	}
	var maxCode int
	for _, customer := range customers {
		maxCode = customer.CustomerCode
	}
	return maxCode
}

func GenCustomerID(customer_code int) string {
	prefix := "PHSPS_CUSTOMER_"
	userID := fmt.Sprintf("%v%d", prefix, customer_code)
	return userID
}
