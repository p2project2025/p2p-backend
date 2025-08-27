package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WithdrawlRequest struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Amount        float64            `bson:"amount" json:"amount"`
	INRRate       float64            `bson:"inr_rate" json:"inr_rate"`
	BankName      string             `bson:"bank_name" json:"bank_name"`
	HolderName    string             `bson:"holder_name" json:"holder_name"`
	AccountNumber string             `bson:"account_number" json:"account_number"`
	IFSCCode      string             `bson:"ifsc_code" json:"ifsc_code"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UserId        primitive.ObjectID `bson:"user_id" json:"user_id"`
	UTR           string             `bson:"utr" json:"utr"`
	ApprovedAt    time.Time          `bson:"approved_at" json:"approved_at"`
	Status        string             `bson:"status" json:"status"`
}

type WithdrawlRes struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Amount        float64            `bson:"amount" json:"amount"`
	INRRate       float64            `bson:"inr_rate" json:"inr_rate"`
	BankName      string             `bson:"bank_name" json:"bank_name"`
	HolderName    string             `bson:"holder_name" json:"holder_name"`
	UTR           string             `bson:"utr" json:"utr"`
	AccountNumber string             `bson:"account_number" json:"account_number"`
	IFSCCode      string             `bson:"ifsc_code" json:"ifsc_code"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	ApprovedAt    time.Time          `bson:"approved_at" json:"approved_at"`
	UserId        primitive.ObjectID `bson:"user_id" json:"user_id"`
	Status        string             `bson:"status" json:"status"`
	User          UserInfo           `bson:"user" json:"user"`
}
