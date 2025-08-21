package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DepositRequest struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Amount          float64            `bson:"amount" json:"amount"`
	TransactionHash string             `bson:"transaction_hash" json:"transaction_hash"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UserId          primitive.ObjectID `bson:"user_id" json:"user_id"`
	Status          string             `bson:"status" json:"status"`
}

type DepositRes struct {
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	Amount          float64            `bson:"amount" json:"amount"`
	TransactionHash string             `bson:"transaction_hash" json:"transaction_hash"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UserId          primitive.ObjectID `bson:"user_id" json:"user_id"`
	Status          string             `bson:"status" json:"status"`
	User            UserInfo           `bson:"user" json:"user"`
}

type UserInfo struct {
	Username string `bson:"username" json:"username"`
	Email    string `bson:"email" json:"email"`
}

type UpdateStatusRequest struct {
	ID     string `json:"id" binding:"required"`
	Status string `json:"status" binding:"required"` // "Approved" or "Rejected"
}
