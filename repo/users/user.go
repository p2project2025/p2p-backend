package users

import (
	"context"
	"errors"
	"fmt"
	"log"
	"p2p/config"
	"p2p/config/db"
	"p2p/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	RegisterUser(user models.User) (primitive.ObjectID, error)
	BlockUser(userID primitive.ObjectID, block bool) error
	CheckPhoneExists(phone string) (bool, error)
	CheckEmailExists(email string) (bool, error)
	GetUserByEmail(email string) (models.User, error)
	GetAllUsers() ([]models.User, error)
}

type UserRepo struct{}

func (r *UserRepo) RegisterUser(user models.User) (primitive.ObjectID, error) {
	collection := db.GetCollection(config.Cfg.DBName, "users")

	// Set auto-generated fields
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()

	// Insert into MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Println(err)
		return primitive.NilObjectID, err
	}

	return user.ID, nil
}

func (r *UserRepo) BlockUser(userID primitive.ObjectID, block bool) error {
	collection := db.GetCollection(config.Cfg.DBName, "users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"is_blocked": block,
		},
	}

	// Update user by ID
	result, err := collection.UpdateOne(ctx, bson.M{"_id": userID}, update)
	if err != nil {
		log.Println("Error blocking user:", err)
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no user found with ID %s", userID.Hex())
	}

	return nil
}

// CheckPhoneExists checks if a phone number is already registered
// Skips check if empty, null, or "0000000000"
func (r *UserRepo) CheckPhoneExists(phone string) (bool, error) {
	if phone == "" || phone == "0000000000" {
		return false, nil
	}

	collection := db.GetCollection(config.Cfg.DBName, "users")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, bson.M{"phone_num": phone})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CheckEmailExists checks if an email is already registered
func (r *UserRepo) CheckEmailExists(email string) (bool, error) {
	if email == "" {
		return false, nil
	}

	collection := db.GetCollection(config.Cfg.DBName, "users")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepo) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	collection := db.GetCollection(config.Cfg.DBName, "users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}

	return user, nil
}

func (r *UserRepo) GetAllUsers() ([]models.User, error) {
	collection := db.GetCollection(config.Cfg.DBName, "users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Filter only users with role = "user"
	filter := bson.M{"role": "user"}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
