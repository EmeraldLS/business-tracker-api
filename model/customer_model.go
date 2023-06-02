package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Customer struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CustomerID    string             `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
	CustomerCode  int                `json:"customer_code,omitempty" bson:"customer_code,omitempty"`
	FirstName     string             `json:"fname,omitempty" bson:"fname,omitempty" validate:"required"`
	LastName      string             `json:"lname,omitempty" bson:"lname,omitempty" validate:"required"`
	BusinessName  string             `json:"bname,omitempty" bson:"bname,omitempty" validate:"required"`
	Phone         int                `json:"phone,omitempty" bson:"phone,omitempty" validate:"required"`
	Email         string             `json:"email,omitempty" bson:"email,omitempty" validate:"required"`
	Location      string             `json:"location,omitempty" bson:"location,omitempty" validate:"required"`
	SetupFee      int                `json:"setup_fee,omitempty" bson:"setup_fee,omitempty" validate:"required"`
	AnnualFee     int                `json:"annual_fee,omitempty" bson:"annual_fee" validate:"required"`
	JobStatusCode int                `json:"job_status_code,omitempty" bson:"job_status_code,omitempty"`
	JobStatus     string             `json:"job_status,omitempty" bson:"job_status,omitempty"`
	JoinDate      string             `json:"join_date,omitempty" bson:"join_date,omitempty"`
	RenewalMonth  string             `bson:"renewal_month,omitempty"`
	RenewalDate   string             `bson:"renewal_date,omitempty"`
	UpdatedAt     string             `bson:"updated_at,omitempty" bson:"updated_at,omitempty"`
	JoinYear      int                `json:"join_year,omitempty" bson:"join_year,omitempty"`
}
