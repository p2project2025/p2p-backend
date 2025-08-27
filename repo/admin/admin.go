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
