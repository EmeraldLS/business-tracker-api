package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName    string             `json:"fname,omitempty" validate:"required" bson:"fname,omitempty"`
	LastName     string             `json:"lname,omitempty" validate:"required" bson:"lname,omitempty"`
	Email        string             `json:"email,omitempty" validate:"required" bson:"email,omitempty"`
	Password     string             `json:"password,omitempty" validate:"required" bson:"password,omitempty"`
	JoinDate     string             `json:"join_date,omitempty" bson:"join_date,omitempty"`
	UpdatedDate  string             `json:"updated_date,omitempty" bson:"updated_date,omitempty"`
	Token        string             `json:"token,omitempty" bson:"token,omitempty"`
	RefreshToken string             `json:"refresh_token,omitempty" bson:"refresh_token,omitempty"`
	ExpiresAt    int64              `json:"expires_at,omitempty" bson:"expires_at,omitempty"`
}

type Login struct {
	Email    string `json:"email,omitempty" validate:"required" bson:"email,omitempty"`
	Password string `json:"password,omitempty" validate:"required" bson:"password,omitempty"`
}
