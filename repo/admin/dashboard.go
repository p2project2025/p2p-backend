package admin

import (
	"context"
	"p2p/config"
	"p2p/config/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type DashboardRepository interface {
	GetCounts() (map[string]int64, error)
}
type DashboardRepo struct{}

func (r *DashboardRepo) GetCounts() (map[string]int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// collections
	userCollection := db.GetCollection(config.Cfg.DBName, "users")
	depositCollection := db.GetCollection(config.Cfg.DBName, "deposit")
	withdrawlCollection := db.GetCollection(config.Cfg.DBName, "withdrawl")

	// queries
	userFilter := bson.M{"is_blocked": false}
	depositFilter := bson.M{"status": "Pending"}
	withdrawlFilter := bson.M{"status": "Pending"}

	// counts
	userCount, err := userCollection.CountDocuments(ctx, userFilter)
	if err != nil {
		return nil, err
	}

	depositCount, err := depositCollection.CountDocuments(ctx, depositFilter)
	if err != nil {
		return nil, err
	}

	withdrawlCount, err := withdrawlCollection.CountDocuments(ctx, withdrawlFilter)
	if err != nil {
		return nil, err
	}

	// return as a map
	return map[string]int64{
		"active_users":       userCount,
		"pending_deposits":   depositCount,
		"pending_withdrawls": withdrawlCount,
	}, nil
}
