package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Sender    primitive.ObjectID `bson:"sender_id" json:"sender_id"`
	Receiver  primitive.ObjectID `bson:"receiver_id" json:"receiver_id"`
	FilesUrl  []string           `bson:"files" json:"files"`
	Message   string             `bson:"message" json:"message"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}

type Chatres struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Sender    UserInfo           `bson:"sender" json:"sender"`
	Receiver  UserInfo           `bson:"receiver" json:"receiver"`
	FilesUrl  []string           `bson:"files" json:"files"`
	Message   string             `bson:"message" json:"message"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}
