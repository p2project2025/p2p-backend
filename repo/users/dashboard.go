package users

import (
	"context"
	"p2p/config"
	"p2p/config/db"
	"p2p/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DashboardRepository interface {
	GetUserDashboard(userID primitive.ObjectID) (*models.UserDash, error)
}

type DashboardRepo struct{}

func (r *DashboardRepo) GetUserDashboard(userID primitive.ObjectID) (*models.UserDash, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userCollection := db.GetCollection(config.Cfg.DBName, "users")
	withdrawlCollection := db.GetCollection(config.Cfg.DBName, "withdrawl")

	// Fetch user balance
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	// Count pending withdrawls for this user
	withdrawlCount, err := withdrawlCollection.CountDocuments(ctx, bson.M{
		"user_id": userID,
		"status":  "Pending",
	})
	if err != nil {
		return nil, err
	}

	return &models.UserDash{
		Balance:            user.Balance,
		PendingWithdrawals: withdrawlCount,
	}, nil
}
