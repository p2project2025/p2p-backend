package deposit

import (
	"context"
	"errors"
	"fmt"
	"log"
	"p2p/config"
	"p2p/config/db"
	"p2p/models"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DepositRepository interface {
	UpdateDepositStatus(depositID string, approve bool) error
	DepositRequest(req models.DepositRequest) error
	GetAll() ([]models.DepositRes, error)
	GetAllByUserID(userID string) ([]models.DepositRes, error)
	GetByID(id string) (*models.DepositRes, error)
	SearchByUsername(username string) ([]models.DepositRes, error)
}

type DepositRepo struct{}

func (r *DepositRepo) UpdateDepositStatus(depositID string, approve bool) error {
	depositCollection := db.GetCollection(config.Cfg.DBName, "deposit")
	userCollection := db.GetCollection(config.Cfg.DBName, "users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := strings.TrimSpace(depositID)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid deposit ID: %w", err)
	}

	// 1️⃣ Find deposit
	var dep models.DepositRequest
	err = depositCollection.FindOne(ctx, bson.M{"_id": oid}).Decode(&dep)
	if err != nil {
		log.Printf("Deposit not found: %v \n object id : %v \n deposit id : %s", err, oid, depositID)
		return fmt.Errorf("deposit not found: %w \n object id : %v \n deposit id : %s", err, oid, depositID)
	}

	// 2️⃣ If already processed, block duplicate updates
	if dep.Status == "Approved" || dep.Status == "Rejected" {
		return fmt.Errorf("deposit already %s", dep.Status)
	}

	if approve {
		// 3️⃣ Approve → update deposit status
		_, err = depositCollection.UpdateOne(
			ctx,
			bson.M{"_id": oid},
			bson.M{"$set": bson.M{"status": "Approved"}},
		)
		if err != nil {
			return fmt.Errorf("failed to approve deposit: %w", err)
		}

		// 4️⃣ Increment user balance
		_, err = userCollection.UpdateOne(
			ctx,
			bson.M{"_id": dep.UserId},
			bson.M{"$inc": bson.M{"balance": dep.Amount}},
		)
		if err != nil {
			return fmt.Errorf("failed to update user balance: %w", err)
		}

	} else {
		// 5️⃣ Reject → only update deposit status
		_, err = depositCollection.UpdateOne(
			ctx,
			bson.M{"_id": oid},
			bson.M{"$set": bson.M{"status": "Rejected"}},
		)
		if err != nil {
			return fmt.Errorf("failed to reject deposit: %w", err)
		}
	}

	return nil
}

func (r *DepositRepo) DepositRequest(req models.DepositRequest) error {
	collection := db.GetCollection(config.Cfg.DBName, "deposit")

	// Set auto-generated fields
	req.ID = primitive.NewObjectID()
	req.CreatedAt = time.Now()

	// Insert into MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, req)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *DepositRepo) GetAllByUserID(userID string) ([]models.DepositRes, error) {
	depositCollection := db.GetCollection(config.Cfg.DBName, "deposit")
	userCollection := db.GetCollection(config.Cfg.DBName, "users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert userID if provided
	var oid primitive.ObjectID
	var err error
	filter := bson.M{}
	if userID != "" {
		oid, err = primitive.ObjectIDFromHex(userID)
		if err != nil {
			log.Println("Invalid user ID:", err)
			return nil, err
		}
		filter["user_id"] = oid
	}

	// 1️⃣ Fetch deposits (no pagination, just sorting)
	findOpts := options.Find().
		SetSort(bson.M{"created_at": -1})

	cursor, err := depositCollection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var deposits []models.DepositRes
	if err := cursor.All(ctx, &deposits); err != nil {
		return nil, err
	}

	// 2️⃣ Fetch user info for each deposit
	var results []models.DepositRes
	for _, dep := range deposits {
		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"_id": dep.UserId}).Decode(&user)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}

		results = append(results, models.DepositRes{
			ID:              dep.ID,
			Amount:          dep.Amount,
			Status:          dep.Status,
			TransactionHash: dep.TransactionHash,
			CreatedAt:       dep.CreatedAt,
			User: models.UserInfo{
				Username: user.Name,
				Email:    user.Email,
			},
		})
	}

	log.Printf("Found %d deposit results", len(results))
	return results, nil
}

func (r *DepositRepo) GetByID(id string) (*models.DepositRes, error) {
	depositCollection := db.GetCollection(config.Cfg.DBName, "deposit")
	userCollection := db.GetCollection(config.Cfg.DBName, "users")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1️⃣ Fetch the deposit
	var dep models.DepositRes
	if err := depositCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&dep); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	// 2️⃣ Fetch the user info
	var user models.User
	if err := userCollection.FindOne(ctx, bson.M{"_id": dep.UserId}).Decode(&user); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	dep.User = models.UserInfo{
		Username: user.Name,
		Email:    user.Email,
	}

	return &dep, nil
}

func (r *DepositRepo) SearchByUsername(username string) ([]models.DepositRes, error) {
	depositCollection := db.GetCollection(config.Cfg.DBName, "deposit")
	userCollection := db.GetCollection(config.Cfg.DBName, "users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1️⃣ Find users by username (case-insensitive search)
	userCursor, err := userCollection.Find(ctx, bson.M{
		"name": bson.M{"$regex": username, "$options": "i"},
	})
	if err != nil {
		return nil, err
	}
	defer userCursor.Close(ctx)

	var users []models.User
	if err := userCursor.All(ctx, &users); err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return []models.DepositRes{}, nil
	}

	// 2️⃣ Extract user IDs
	var userIDs []primitive.ObjectID
	for _, u := range users {
		userIDs = append(userIDs, u.ID)
	}

	// 3️⃣ Find deposits for these user IDs (no pagination, just sorting)
	findOpts := options.Find().SetSort(bson.M{"created_at": -1})

	depCursor, err := depositCollection.Find(ctx, bson.M{"user_id": bson.M{"$in": userIDs}}, findOpts)
	if err != nil {
		return nil, err
	}
	defer depCursor.Close(ctx)

	var deposits []models.DepositRes
	if err := depCursor.All(ctx, &deposits); err != nil {
		return nil, err
	}

	// 4️⃣ Map user info
	userMap := make(map[primitive.ObjectID]models.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	var results []models.DepositRes
	for _, dep := range deposits {
		if user, ok := userMap[dep.UserId]; ok {
			dep.User = models.UserInfo{
				Username: user.Name,
				Email:    user.Email,
			}
		}
		results = append(results, dep)
	}

	return results, nil
}

func (r *DepositRepo) GetAll() ([]models.DepositRes, error) {
	depositCollection := db.GetCollection(config.Cfg.DBName, "deposit")
	userCollection := db.GetCollection(config.Cfg.DBName, "users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1️⃣ Get deposits (no pagination, only sort)
	findOpts := options.Find().
		SetSort(bson.M{"created_at": -1})

	depCursor, err := depositCollection.Find(ctx, bson.M{}, findOpts)
	if err != nil {
		return nil, err
	}
	defer depCursor.Close(ctx)

	var deposits []models.DepositRes
	if err := depCursor.All(ctx, &deposits); err != nil {
		return nil, err
	}

	// 2️⃣ Collect user IDs
	userIDs := make([]primitive.ObjectID, 0, len(deposits))
	for _, dep := range deposits {
		userIDs = append(userIDs, dep.UserId)
	}

	// 3️⃣ Fetch all users in one query
	userCursor, err := userCollection.Find(ctx, bson.M{"_id": bson.M{"$in": userIDs}})
	if err != nil {
		return nil, err
	}
	defer userCursor.Close(ctx)

	var users []models.User
	if err := userCursor.All(ctx, &users); err != nil {
		return nil, err
	}

	// 4️⃣ Map user IDs to users
	userMap := make(map[primitive.ObjectID]models.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	// 5️⃣ Merge results
	var results []models.DepositRes
	for _, dep := range deposits {
		if user, ok := userMap[dep.UserId]; ok {
			dep.User = models.UserInfo{
				Username: user.Name,
				Email:    user.Email,
			}
		}
		results = append(results, dep)
	}

	return results, nil
}

// -------------------- PIPELINE HELPERS --------------------
// Use consistent BSON type (bson.D) throughout
func lookupStage() bson.D {
	return bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "user_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "user"},
		}},
	}
}

func unwindUser() bson.D {
	return bson.D{
		{Key: "$unwind", Value: "$user"},
	}
}

func projectFields() bson.D {
	return bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "amount", Value: 1},
			{Key: "created_at", Value: 1},
			{Key: "user", Value: bson.D{
				{Key: "username", Value: 1},
				{Key: "email", Value: 1},
			}},
		}},
	}
}

func mongoPipeline(skip, limit int64, username string, userID *primitive.ObjectID) mongo.Pipeline {
	var pipeline mongo.Pipeline

	if userID != nil {
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "user_id", Value: *userID},
			}},
		})
	}

	pipeline = append(pipeline, lookupStage(), unwindUser())

	if username != "" {
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "user.username", Value: bson.D{
					{Key: "$regex", Value: username},
					{Key: "$options", Value: "i"},
				}},
			}},
		})
	}

	pipeline = append(pipeline,
		projectFields(),
		bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		bson.D{{Key: "$skip", Value: skip}},
		bson.D{{Key: "$limit", Value: limit}},
	)

	return pipeline
}

// Alternative: If you need all fields from your original struct, update the projection
func projectAllFields() bson.D {
	return bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "amount", Value: 1},
			{Key: "transaction_hash", Value: 1},
			{Key: "created_at", Value: 1},
			{Key: "user_id", Value: 1},
			{Key: "status", Value: 1},
			{Key: "user_name", Value: "$user.username"}, // Map username from user object
			{Key: "user", Value: bson.D{
				{Key: "username", Value: 1},
				{Key: "email", Value: 1},
			}},
		}},
	}
}
