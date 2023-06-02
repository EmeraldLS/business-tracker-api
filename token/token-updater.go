package token

import (
	"github.com/EmeraldLS/phsps-api/config"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateToken(signedToken string, refreshToken string, updateTime string, expirationTime int64, email string) (int64, error) {
	var updateObj = bson.M{}
	updateObj["token"] = signedToken
	updateObj["refresh_token"] = refreshToken
	updateObj["updated_date"] = updateTime
	updateObj["expires_at"] = expirationTime
	update := bson.M{"$set": updateObj}
	count, err := config.UpdateToken(email, update)
	if err != nil {
		return 0, err
	}
	return count, nil
}
