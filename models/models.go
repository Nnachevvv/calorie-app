package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username,omitempty"`
	Password string             `json:"password,omitempty"`
}

type RegisterUser struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username        string             `json:"username,omitempty"`
	Password        string             `json:"password,omitempty"`
	ConfirmPassword string             `json:"confirmpassword,omitempty"`
	Email           string             `json:"email,omitempty"`
}
