package config

import (
	"context"
	"errors"
	"time"

	"github.com/EmeraldLS/phsps-api/db"
	"github.com/EmeraldLS/phsps-api/model"
	"github.com/badoux/checkmail"
	"github.com/golang-module/carbon"
	"go.mongodb.org/mongo-driver/bson"
)

var UserCollection = db.UserCollection

func Register(user model.User) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	_, err := UserCollection.InsertOne(ctx, user)
	defer cancel()
	if err != nil {
		return err
	}
	return nil
}

func CheckEmailExist(email string) error {
	filter := bson.M{"email": email}
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	count, err := UserCollection.CountDocuments(ctx, filter)
	defer cancel()
	if err != nil {
		return err
	}
	if count > 1 {
		return errors.New("email already exist")
	}
	return nil
}

func ValidateEmail(email string) error {
	if err := checkmail.ValidateFormat(email); err != nil {
		return errors.New("invalid email format. Please provide a valid email")
	}
	return nil
}

func Login(user model.Login) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	filter := bson.M{"email": user.Email}
	var foundUser model.User
	err := UserCollection.FindOne(ctx, filter).Decode(&foundUser)
	defer cancel()
	if err != nil {
		return model.User{}, errors.New("user not found")
	}
	return foundUser, nil
}

func Logout(token string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	var updateObj = bson.M{}
	updateObj["expires_at"] = time.Now().Unix()
	updateObj["updated_date"] = carbon.Now().ToDateTimeString()
	update := bson.M{"$set": updateObj}
	result, err := UserCollection.UpdateOne(ctx, bson.M{"token": token}, update)
	defer cancel()
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

func UpdateToken(email string, updateObj bson.M) (int64, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	result, err := UserCollection.UpdateOne(ctx, bson.M{"email": email}, updateObj)
	defer cancel()
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

func FindUserByToken(token string) (model.User, error) {
	var user model.User
	filter := bson.M{"token": token}
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	err := UserCollection.FindOne(ctx, filter).Decode(&user)
	defer cancel()
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func FindUserByRefreshToken(refresh_token string) (model.User, error) {
	var user model.User
	filter := bson.M{"refresh_token": refresh_token}
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	err := UserCollection.FindOne(ctx, filter).Decode(&user)
	defer cancel()
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
