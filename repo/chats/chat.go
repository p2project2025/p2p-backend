package chats

import (
	"context"
	"p2p/config"
	"p2p/config/db"
	"p2p/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ChatRepoInterface interface {
	CreateChat(chat *models.Chat) error
	GetChatsBetweenUsers(userA, userB primitive.ObjectID) ([]models.Chatres, error)
	GetUniqueChatUsers(userID primitive.ObjectID) ([]models.User, error)
}

type ChatRepo struct{}

// CreateChat inserts a new chat into the "chats" collection
func (r *ChatRepo) CreateChat(chat *models.Chat) error {
	collection := db.GetCollection(config.Cfg.DBName, "chats")

	chat.ID = primitive.NewObjectID()
	chat.Timestamp = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, chat)
	if err != nil {
		return err
	}
	return nil
}

func (r *ChatRepo) GetChatsBetweenUsers(userA, userB primitive.ObjectID) ([]models.Chatres, error) {
	chatCollection := db.GetCollection(config.Cfg.DBName, "chats")
	userCollection := db.GetCollection(config.Cfg.DBName, "users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1️⃣ Fetch chats
	filter := bson.M{
		"$or": []bson.M{
			{"sender_id": userA, "receiver_id": userB},
			{"sender_id": userB, "receiver_id": userA},
		},
	}
	opts := options.Find().SetSort(bson.M{"timestamp": 1})

	cursor, err := chatCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chats []models.Chat
	if err := cursor.All(ctx, &chats); err != nil {
		return nil, err
	}

	// 2️⃣ Build response with user info
	var results []models.Chatres
	for _, chat := range chats {
		var sender, receiver models.User

		_ = userCollection.FindOne(ctx, bson.M{"_id": chat.Sender}).Decode(&sender)
		_ = userCollection.FindOne(ctx, bson.M{"_id": chat.Receiver}).Decode(&receiver)

		results = append(results, models.Chatres{
			ID:        chat.ID,
			Sender:    models.UserInfo{Username: sender.Name, Email: sender.Email},
			Receiver:  models.UserInfo{Username: receiver.Name, Email: receiver.Email},
			FilesUrl:  chat.FilesUrl,
			Message:   chat.Message,
			Timestamp: chat.Timestamp,
		})
	}

	return results, nil
}

func (r *ChatRepo) GetUniqueChatUsers(userID primitive.ObjectID) ([]models.User, error) {
	chatCollection := db.GetCollection(config.Cfg.DBName, "chats")
	userCollection := db.GetCollection(config.Cfg.DBName, "users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Fetch all chats involving this user
	filter := bson.M{
		"$or": []bson.M{
			{"sender_id": userID},
			{"receiver_id": userID},
		},
	}

	cursor, err := chatCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 2. Collect unique other user IDs
	uniqueIDs := make(map[primitive.ObjectID]struct{})
	for cursor.Next(ctx) {
		var chat models.Chat
		if err := cursor.Decode(&chat); err != nil {
			return nil, err
		}

		if chat.Sender == userID {
			uniqueIDs[chat.Receiver] = struct{}{}
		} else {
			uniqueIDs[chat.Sender] = struct{}{}
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	// Convert map keys to slice
	var ids []primitive.ObjectID
	for id := range uniqueIDs {
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return []models.User{}, nil
	}

	// 3. Fetch users by IDs
	userFilter := bson.M{"_id": bson.M{"$in": ids}}
	userCursor, err := userCollection.Find(ctx, userFilter)
	if err != nil {
		return nil, err
	}
	defer userCursor.Close(ctx)

	var users []models.User
	if err := userCursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
