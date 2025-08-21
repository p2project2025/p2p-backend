package common

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func LookupStage() bson.D {
	return bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "users",
		"localField":   "user_id",
		"foreignField": "_id",
		"as":           "user",
	}}}
}

func UnwindUser() bson.D {
	return bson.D{{Key: "$unwind", Value: bson.M{
		"path":                       "$user",
		"preserveNullAndEmptyArrays": true,
	}}}
}

func ProjectFields() bson.D {
	return bson.D{{Key: "$project", Value: bson.M{
		"id":               "$_id",
		"user_name":        "$user.username",
		"amount":           1,
		"transaction_hash": 1,
		"created_at":       1,
		"user_id":          1,
		"status":           1,
	}}}
}

func MongoPipeline(skip, limit int64, username string, userID *primitive.ObjectID) mongo.Pipeline {
	var pipeline mongo.Pipeline

	if userID != nil {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.M{"user_id": *userID}}})
	}

	pipeline = append(pipeline, LookupStage(), UnwindUser())

	if username != "" {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.M{
			"user.username": bson.M{"$regex": username, "$options": "i"},
		}}})
	}

	pipeline = append(pipeline,
		ProjectFields(),
		bson.D{{Key: "$sort", Value: bson.M{"created_at": -1}}},
		bson.D{{Key: "$skip", Value: skip}},
		bson.D{{Key: "$limit", Value: limit}},
	)

	return pipeline
}
