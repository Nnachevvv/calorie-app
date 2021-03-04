package middleware

import (
	"context"
	"errors"

	"github.com/Nnachevv/calorieapp/models"
	"go.mongodb.org/mongo-driver/bson"
)

type MongoDatabase interface {
	Find(string) (bson.M, error)
	Add(models.RegisterUser) error
}

type Service struct {
}

//ErrUserIsNotFound error represent when user is not present in database
var ErrUserIsNotFound = errors.New("this user is not found")

// Find gets data if exist from mongo db client
func (s *Service) Find(username string) (bson.M, error) {
	var user bson.M
	collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if user == nil {
		return nil, ErrUserIsNotFound
	}

	return user, nil
}

// Add user to database
func (s *Service) Add(user models.RegisterUser) error {
	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}

	return nil
}
