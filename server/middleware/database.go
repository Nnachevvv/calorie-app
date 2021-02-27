package middleware

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDatabase interface {
	Find(username string, collection *mongo.Collection) (bson.M, error)
}

type Service struct {
}

// Find gets data if exist from mongo db client
func (s *Service) Find(username string, collection *mongo.Collection) (bson.M, error) {
	var user bson.M
	collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if user == nil {
		return nil, errors.New("failed to find this account")
	}

	return user, nil
}
