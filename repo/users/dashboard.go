package users

import (
	"context"
	"errors"
	"p2p/config"
	"p2p/config/db"
	"p2p/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	// 1️⃣ Fetch user balance (handle no user case)
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// ✅ No user found → return empty dashboard
			return &models.UserDash{
				Balance:            0,
				PendingWithdrawals: 0,
			}, nil
		}
		return nil, err
	}

	// 2️⃣ Count pending withdrawals (safe even if no docs)
	withdrawlCount, err := withdrawlCollection.CountDocuments(ctx, bson.M{
		"user_id": userID,
		"status":  "Pending",
	})
	if err != nil {
		return nil, err
	}

	// 3️⃣ Return dashboard data
	return &models.UserDash{
		Balance:            user.Balance,
		PendingWithdrawals: withdrawlCount,
	}, nil
}
