package admin

import (
	"context"
	"log"
	"p2p/config"
	"p2p/config/db"
	"p2p/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AdminRepository interface {
	Upsert(admin models.AdminConfigData) (primitive.ObjectID, error)
	Fetch() (*models.AdminConfigData, error)
	GetLedgerStats() (*models.LedgerRes, error)
}


type AdminRepo struct{}

func (r *AdminRepo) Upsert(admin models.AdminConfigData) (primitive.ObjectID, error) {
	collection := db.GetCollection(config.Cfg.DBName, "admin")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Build update fields dynamically only if not empty
	updateFields := bson.M{}
	if admin.SecureWalletAddress != "" {
		updateFields["secure_wallet_address"] = admin.SecureWalletAddress
	}
	if admin.USDTRate != "" {
		updateFields["usdt_rate"] = admin.USDTRate
	}
	if admin.QRCodeURL != "" {
		updateFields["qr_code_url"] = admin.QRCodeURL
	}

	// If no fields to update, return early
	if len(updateFields) == 0 {
		return primitive.NilObjectID, nil
	}

	// Upsert option (✅ v2 correct usage)
	opts := options.UpdateOne().SetUpsert(true)

	// Filter - assuming only one admin config exists
	filter := bson.M{}

	// Update operation
	update := bson.M{"$set": updateFields}

	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return primitive.NilObjectID, err
	}

	// If new document inserted, return the inserted ID
	if result.UpsertedID != nil {
		if oid, ok := result.UpsertedID.(primitive.ObjectID); ok {
			return oid, nil
		}
	}

	// Otherwise, fetch the existing document’s ID
	var updatedDoc bson.M
	err = collection.FindOne(ctx, filter).Decode(&updatedDoc)
	if err != nil {
		return primitive.NilObjectID, err
	}

	if id, ok := updatedDoc["_id"].(primitive.ObjectID); ok {
		return id, nil
	}

	return primitive.NilObjectID, nil
}

func (r *AdminRepo) Fetch() (*models.AdminConfigData, error) {
	collection := db.GetCollection(config.Cfg.DBName, "admin")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var admin models.AdminConfigData
	err := collection.FindOne(ctx, bson.M{}).Decode(&admin)
	if err != nil {
		// Handle "no documents" gracefully
		if err == mongo.ErrNoDocuments {
			return &models.AdminConfigData{}, nil
		}
		// Any other DB error
		log.Println("failed to fetch cnf data :", err)
		return nil, err
	}

	return &admin, nil
}
func (r *AdminRepo) GetLedgerStats() (*models.LedgerRes, error) {
	depositCollection := db.GetCollection(config.Cfg.DBName, "deposit")
	withdrawCollection := db.GetCollection(config.Cfg.DBName, "withdrawl")
	userCollection := db.GetCollection(config.Cfg.DBName, "users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ledger := &models.LedgerRes{}

	// 1️⃣ Total deposits
	depCursor, err := depositCollection.Aggregate(ctx, mongo.Pipeline{
		{{Key: "$group", Value: bson.M{"_id": nil, "total": bson.M{"$sum": "$amount"}}}},
	})
	if err == nil && depCursor.Next(ctx) {
		var res struct {
			Total float64 `bson:"total"`
		}
		if err := depCursor.Decode(&res); err == nil {
			ledger.TotalDeposits = res.Total
		}
	}
	depCursor.Close(ctx)

	// 2️⃣ Withdrawals by status
	withCursor, err := withdrawCollection.Aggregate(ctx, mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id":   "$status",
			"total": bson.M{"$sum": "$amount"},
			"count": bson.M{"$sum": 1},
		}}},
	})
	if err == nil {
		for withCursor.Next(ctx) {
			var res struct {
				Status string  `bson:"_id"`
				Total  float64 `bson:"total"`
				Count  int64   `bson:"count"`
			}
			if err := withCursor.Decode(&res); err == nil {
				switch res.Status {
				case "Approved":
					ledger.TotalWithdrawals += res.Total
				case "Pending":
					ledger.TotalPendingWithdrawals = res.Count
					ledger.PendingWithdrawalsTotal = res.Total
				case "Rejected":
					ledger.RejectedWithdrawalsTotal = res.Total
				}
			}
		}
	}
	withCursor.Close(ctx)

	// 3️⃣ Current total balance (sum of all user balances)
	userCursor, err := userCollection.Aggregate(ctx, mongo.Pipeline{
		{{Key: "$group", Value: bson.M{"_id": nil, "total": bson.M{"$sum": "$balance"}}}},
	})
	if err == nil && userCursor.Next(ctx) {
		var res struct {
			Total float64 `bson:"total"`
		}
		if err := userCursor.Decode(&res); err == nil {
			ledger.CurrentTotalBalance = res.Total
		}
	}
	userCursor.Close(ctx)

	// 4️⃣ Today’s stats
	startOfDay := time.Now().Truncate(24 * time.Hour)

	// ➤ Deposits today by status
	depTodayCursor, _ := depositCollection.Aggregate(ctx, mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"created_at": bson.M{"$gte": startOfDay}}}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$status",
			"total": bson.M{"$sum": "$amount"},
		}}},
	})
	for depTodayCursor.Next(ctx) {
		var res struct {
			Status string  `bson:"_id"`
			Total  float64 `bson:"total"`
		}
		_ = depTodayCursor.Decode(&res)
		switch res.Status {
		case "Approved":
			ledger.TodayStats.TotalDepositsApproved = res.Total
			ledger.TodayStats.TotalDeposits += res.Total
		case "Pending":
			ledger.TodayStats.TotalDepositsPending = res.Total
			ledger.TodayStats.TotalDeposits += res.Total
		}
	}
	depTodayCursor.Close(ctx)

	// ➤ Withdrawals today by status
	withTodayCursor, _ := withdrawCollection.Aggregate(ctx, mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"created_at": bson.M{"$gte": startOfDay}}}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$status",
			"total": bson.M{"$sum": "$amount"},
		}}},
	})
	for withTodayCursor.Next(ctx) {
		var res struct {
			Status string  `bson:"_id"`
			Total  float64 `bson:"total"`
		}
		_ = withTodayCursor.Decode(&res)
		switch res.Status {
		case "Approved":
			ledger.TodayStats.TotalWithdrawalsApproved = res.Total
			ledger.TodayStats.TotalWithdrawals += res.Total
		case "Pending":
			ledger.TodayStats.TotalWithdrawalsPending = res.Total
			ledger.TodayStats.TotalWithdrawals += res.Total
		}
	}
	withTodayCursor.Close(ctx)

	// ➤ New users today
	newUsersCount, _ := userCollection.CountDocuments(ctx, bson.M{"created_at": bson.M{"$gte": startOfDay}})
	ledger.TodayStats.NewUsers = newUsersCount

	return ledger, nil
}