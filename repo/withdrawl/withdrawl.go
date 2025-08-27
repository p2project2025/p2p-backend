package withdrawl

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
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type WithdrawlRepository interface {
	WithdrawlRequest(req models.WithdrawlRequest) error
	GetAll() ([]models.WithdrawlRes, error)
	GetByID(id string) (*models.WithdrawlRes, error)
	SearchByUsername(username string) ([]models.WithdrawlRes, error)
	GetAllByUserID(userID string) ([]models.WithdrawlRes, error)
	UpdateWithdrawStatus(withdrawID string, approve bool, utr string) error
}

type WithdrawlRepo struct{}

// -------------------- CREATE WITHDRAWL --------------------
func (r *WithdrawlRepo) WithdrawlRequest(req models.WithdrawlRequest) error {
	collection := db.GetCollection(config.Cfg.DBName, "withdrawl")

	req.ID = primitive.NewObjectID()
	req.CreatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userCollection := db.GetCollection(config.Cfg.DBName, "users")

	// 3️⃣ Get user balance
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"_id": req.UserId}).Decode(&user)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// 4️⃣ Check balance
	if user.Balance < req.Amount {
		return fmt.Errorf("insufficient balance for withdrawal")
	}

	// 6️⃣ Deduct user balance
	_, err = userCollection.UpdateOne(
		ctx,
		bson.M{"_id": req.UserId},
		bson.M{"$inc": bson.M{"balance": -req.Amount}})
	if err != nil {
		return fmt.Errorf("failed to update user balance: %w", err)
	}
	_, err = collection.InsertOne(ctx, req)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *WithdrawlRepo) UpdateWithdrawStatus(withdrawID string, approve bool, utr string) error {
	withdrawCollection := db.GetCollection(config.Cfg.DBName, "withdrawl")
	userCollection := db.GetCollection(config.Cfg.DBName, "users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert withdrawal ID
	oid, err := primitive.ObjectIDFromHex(withdrawID)
	if err != nil {
		return fmt.Errorf("invalid withdraw ID: %w", err)
	}

	// 1️⃣ Find withdrawal
	var wd models.WithdrawlRequest
	err = withdrawCollection.FindOne(ctx, bson.M{"_id": oid}).Decode(&wd)
	if err != nil {
		log.Printf("Withdraw not found: %v \n object id : %v \n withdraw id : %s", err, oid, withdrawID)
		return fmt.Errorf("withdraw not found: %w", err)
	}

	// 2️⃣ Prevent duplicate approvals/rejections
	if wd.Status == "Approved" || wd.Status == "Rejected" {
		return fmt.Errorf("withdraw already %s", wd.Status)
	}

	if approve {

		// 5️⃣ Approve withdrawal + Update UTR
		_, err = withdrawCollection.UpdateOne(
			ctx,
			bson.M{"_id": oid},
			bson.M{"$set": bson.M{
				"status":      "Approved",
				"utr":         utr,
				"approved_at": time.Now(),
			}},
		)
		if err != nil {
			return fmt.Errorf("failed to approve withdraw: %w", err)
		}

	} else {

		_, err = userCollection.UpdateOne(
			ctx,
			bson.M{"_id": wd.UserId},
			bson.M{"$inc": bson.M{"balance": +wd.Amount}},
		)
		if err != nil {
			return fmt.Errorf("failed to update user balance: %w", err)
		}

		// 7️⃣ Reject withdrawal
		_, err = withdrawCollection.UpdateOne(
			ctx,
			bson.M{"_id": oid},
			bson.M{"$set": bson.M{"status": "Rejected"}},
		)
		if err != nil {
			return fmt.Errorf("failed to reject withdraw: %w", err)
		}
	}

	return nil
}

func (r *WithdrawlRepo) GetAll() ([]models.WithdrawlRes, error) {
	withdrawlCollection := db.GetCollection(config.Cfg.DBName, "withdrawl")
	userCollection := db.GetCollection(config.Cfg.DBName, "users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1️⃣ Fetch withdrawals (no pagination, only sort)
	findOpts := options.Find().SetSort(bson.M{"created_at": -1})
	withCursor, err := withdrawlCollection.Find(ctx, bson.M{}, findOpts)
	if err != nil {
		return nil, err
	}
	defer withCursor.Close(ctx)

	var withdrawls []models.WithdrawlRes
	if err := withCursor.All(ctx, &withdrawls); err != nil {
		return nil, err
	}

	// ✅ Handle empty withdrawals (return empty slice, not nil, no error)
	if len(withdrawls) == 0 {
		return []models.WithdrawlRes{}, nil
	}

	// 2️⃣ Collect user IDs
	userIDs := make([]primitive.ObjectID, 0, len(withdrawls))
	for _, w := range withdrawls {
		userIDs = append(userIDs, w.UserId)
	}

	// 3️⃣ Fetch all users in one go
	userCursor, err := userCollection.Find(ctx, bson.M{"_id": bson.M{"$in": userIDs}})
	if err != nil {
		return nil, err
	}
	defer userCursor.Close(ctx)

	var users []models.User
	if err := userCursor.All(ctx, &users); err != nil {
		return nil, err
	}

	// 4️⃣ Map user ID to user
	userMap := make(map[primitive.ObjectID]models.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	// 5️⃣ Merge user info into withdrawals
	for i, w := range withdrawls {
		if u, ok := userMap[w.UserId]; ok {
			withdrawls[i].User = models.UserInfo{
				Username: u.Name,
				Email:    u.Email,
			}
		}
	}

	return withdrawls, nil
}

func (r *WithdrawlRepo) GetByID(id string) (*models.WithdrawlRes, error) {
	withdrawlCollection := db.GetCollection(config.Cfg.DBName, "withdrawl")
	userCollection := db.GetCollection(config.Cfg.DBName, "users")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1️⃣ Fetch withdrawal
	var withdrawl models.WithdrawlRes
	if err := withdrawlCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&withdrawl); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	// 2️⃣ Fetch user
	var user models.User
	if err := userCollection.FindOne(ctx, bson.M{"_id": withdrawl.UserId}).Decode(&user); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	withdrawl.User = models.UserInfo{
		Username: user.Name,
		Email:    user.Email,
	}

	return &withdrawl, nil
}

func (r *WithdrawlRepo) GetAllByUserID(userID string) ([]models.WithdrawlRes, error) {
	withdrawlCollection := db.GetCollection(config.Cfg.DBName, "withdrawl")
	userCollection := db.GetCollection(config.Cfg.DBName, "users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// 1️⃣ Fetch withdrawals for the given user
	findOpts := options.Find().SetSort(bson.M{"created_at": -1})
	withCursor, err := withdrawlCollection.Find(ctx, bson.M{"user_id": oid}, findOpts)
	if err != nil {
		return nil, err
	}
	defer withCursor.Close(ctx)

	var withdrawls []models.WithdrawlRes
	if err := withCursor.All(ctx, &withdrawls); err != nil {
		return nil, err
	}

	// ✅ Handle empty withdrawals (return empty slice, no error)
	if len(withdrawls) == 0 {
		return []models.WithdrawlRes{}, nil
	}

	// 2️⃣ Fetch user info (ignore if not found)
	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	// 3️⃣ Merge user info
	for i := range withdrawls {
		withdrawls[i].User = models.UserInfo{
			Username: user.Name,
			Email:    user.Email,
		}
	}

	return withdrawls, nil
}

func (r *WithdrawlRepo) SearchByUsername(username string) ([]models.WithdrawlRes, error) {
	withdrawlCollection := db.GetCollection(config.Cfg.DBName, "withdrawl")
	userCollection := db.GetCollection(config.Cfg.DBName, "users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1️⃣ Find matching users
	userCursor, err := userCollection.Find(ctx, bson.M{"name": bson.M{"$regex": username, "$options": "i"}})
	if err != nil {
		return nil, err
	}
	defer userCursor.Close(ctx)

	var users []models.User
	if err := userCursor.All(ctx, &users); err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return []models.WithdrawlRes{}, nil
	}

	// 2️⃣ Extract user IDs
	var userIDs []primitive.ObjectID
	for _, u := range users {
		userIDs = append(userIDs, u.ID)
	}

	// 3️⃣ Fetch withdrawals for those users (no pagination)
	findOpts := options.Find().
		SetSort(bson.M{"created_at": -1})

	withCursor, err := withdrawlCollection.Find(ctx, bson.M{"user_id": bson.M{"$in": userIDs}}, findOpts)
	if err != nil {
		return nil, err
	}
	defer withCursor.Close(ctx)

	var withdrawls []models.WithdrawlRes
	if err := withCursor.All(ctx, &withdrawls); err != nil {
		return nil, err
	}

	// 4️⃣ Map user info
	userMap := make(map[primitive.ObjectID]models.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	for i, w := range withdrawls {
		if u, ok := userMap[w.UserId]; ok {
			withdrawls[i].User = models.UserInfo{
				Username: u.Name,
				Email:    u.Email,
			}
		}
	}

	return withdrawls, nil
}
