package config

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/EmeraldLS/phsps-api/db"
	"github.com/EmeraldLS/phsps-api/model"
	"github.com/golang-module/carbon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var CustomerCollection = db.CustomerCollection

func InsertCustomer(customer model.Customer) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	_, err := CustomerCollection.InsertOne(ctx, customer)
	defer cancel()
	if err != nil {
		return err
	}
	return nil
}

func GetAllCustomter() ([]model.Customer, error) {
	var customers []model.Customer
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	filter := bson.M{}
	cursor, err := CustomerCollection.Find(ctx, filter)
	defer cancel()
	if err != nil {
		return []model.Customer{}, err
	}

	for cursor.Next(ctx) {

		var customer model.Customer
		cursor.Decode(&customer)
		customers = append(customers, customer)
	}
	defer cursor.Close(ctx)
	return customers, nil
}

func GetCustomerByID(customerCode int) (model.Customer, error) {
	var customer model.Customer
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	filter := bson.M{"customer_code": customerCode}
	err := CustomerCollection.FindOne(ctx, filter).Decode(&customer)
	defer cancel()
	if err != nil {
		return model.Customer{}, errors.New("customer with provided id is not found")
	}
	return customer, nil
}

func SearchCustomerByBusinessName(queryParam string) ([]model.Customer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var customers []model.Customer
	var filter = bson.M{}
	if queryParam != "" {
		filter = bson.M{
			"bname": bson.M{
				"$regex": primitive.Regex{
					Pattern: queryParam,
					Options: "i",
				},
			},
		}
	}
	cursor, err := CustomerCollection.Find(ctx, filter)
	defer cancel()
	if err != nil {
		return []model.Customer{}, err
	}
	for cursor.Next(ctx) {
		var customer model.Customer
		cursor.Decode(&customer)
		customers = append(customers, customer)
	}
	defer cursor.Close(ctx)
	return customers, nil

}
func UpdateJobStatus(statusCode int, customerCode int) (int64, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	var customer model.Customer
	filter := bson.M{"customer_code": customerCode}
	err := CustomerCollection.FindOne(ctx, filter).Decode(&customer)
	defer cancel()
	if err != nil {
		return 0, errors.New("no customer with provided info")
	}
	updateObj := bson.M{}
	if statusCode == 1 {
		updateObj["job_status_code"] = statusCode
		updateObj["job_status"] = "Job started and in progress."
	} else if statusCode == 2 {
		updateObj["job_status_code"] = statusCode
		updateObj["job_status"] = "Job completed."
	}
	updateObj["updated_at"] = carbon.Now().ToDateTimeString()
	update := bson.M{"$set": updateObj}
	result, err := CustomerCollection.UpdateOne(ctx, filter, update)
	defer cancel()
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

func DeleteCustomer(customerCode int) (int64, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	filter := bson.M{"customer_code": customerCode}
	result, err := CustomerCollection.DeleteOne(ctx, filter)
	defer cancel()
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func GetAllCustomerInASpecificMonth(month string) (customers []model.Customer, err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	filter := bson.M{"renewal_month": month}
	cursor, err := CustomerCollection.Find(ctx, filter)
	defer cancel()
	if err != nil {
		return []model.Customer{}, err
	}
	for cursor.Next(ctx) {
		var customer model.Customer
		cursor.Decode(&customer)
		customers = append(customers, customer)
	}
	defer cursor.Close(ctx)
	return customers, nil
}

func GetTotalAnnualSubFeeForAYear(year int) (string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	filter := bson.M{"join_year": year}
	var customers []model.Customer
	cursor, err := CustomerCollection.Find(ctx, filter)
	defer cancel()
	if err != nil {
		return "", err
	}
	for cursor.Next(ctx) {
		var customer model.Customer
		cursor.Decode(&customer)
		customers = append(customers, customer)
	}
	var total int
	for _, customer := range customers {
		total += customer.AnnualFee
	}
	return fmt.Sprintf("Total annual fee for year -> %v is: %v", year, total), nil
}

func GetTotalAnnualSubFeeOfAMonth(month string) (string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	filter := bson.M{"renewal_month": month}
	var customers []model.Customer
	cursor, err := CustomerCollection.Find(ctx, filter)
	defer cancel()
	if err != nil {
		return "", err
	}
	for cursor.Next(ctx) {
		var customer model.Customer
		cursor.Decode(&customer)
		customers = append(customers, customer)
	}
	var total int
	for _, customer := range customers {
		total += customer.AnnualFee
	}
	return fmt.Sprintf("Total Subscription fee for Month -> %v is: %v", month, total), nil
}
