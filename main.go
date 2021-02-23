package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	var day, year, month, calories string
	fmt.Println("Enter day:")
	fmt.Scanln(&day)

	fmt.Println("Enter month:")
	fmt.Scanln(&month)
	fmt.Println("Enter year:")
	fmt.Scanln(&year)

	fmt.Println("Enter calories:")

	fmt.Scanln(&calories)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Fatal err !")
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	collection := client.Database("testing").Collection("calories")
	collection.InsertOne(ctx, primitive.M{day + "/" + month + "/" + year: calories})
}
