package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Nnachevv/calorieapp/models"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB connection string
// for localhost mongoDB
const connectionString = "mongodb://localhost:27017"

// Database Name
const dbName = "calories-app"

// Collection name
const collName = "users"

// collection object/instance
var collection *mongo.Collection

var MongoService MongoDatabase

// create connection with mongo db
func Connect() {

	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection = client.Database(dbName).Collection(collName)

	fmt.Println("Collection instance created!")
}

//ErrWrongPassword error represent when password is wrong
var ErrWrongPassword = errors.New("this password for this username is wrong")

// LoginToSystem logins user into the system
func LoginToSystem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var userRequest models.User
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		log.Fatal(err)
	}
	err = verifyUser(userRequest)
	if err == ErrUserIsNotFound {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(ErrUserIsNotFound.Error()))
	} else if err == ErrWrongPassword {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(ErrWrongPassword.Error()))
	}
}

// Check if user is in database
func verifyUser(login models.User) error {
	user, err := MongoService.Find(login.Username, collection)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(login.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	// Comparing the password with the hash
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(user["password"].(string)))
	if err != nil {
		return ErrWrongPassword
	}

	return nil
}
