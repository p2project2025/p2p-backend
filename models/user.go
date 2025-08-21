package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	PhoneNum  string             `bson:"phone_num" json:"phone_num"`
	Password  string             `bson:"password" json:"password"`
	Role      string             `bson:"role" json:"role"`
	Balance   float64            `bson:"balance" json:"balance"`
	IsBlocked bool               `bson:"is_blocked" json:"is_blocked"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserDash struct {
	Balance            float64 `json:"balance"`
	PendingWithdrawals int64   `json:"pending_withdrawals"`
	SellPrice          string  `json:"sell_price"`
	WalletAddress      string  `json:"wallet_address"` //..
}
