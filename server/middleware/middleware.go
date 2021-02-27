package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Nnachevv/calorieapp/models"

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
func init() {

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

	verifyUser(userRequest)
	json.NewEncoder(w).Encode(userRequest)
}

// Check if user is in database
func verifyUser(login models.User) {

	//vaultPwd := argon2.IDKey([]byte(login.Username), []byte(login.Password), 1, 64*1024, 4, 32)

	user, err := MongoService.Find(login.Username, collection)
	if err != nil {
		fmt.Println("TODO")
	}

	fmt.Println(user)

	//TODO check this try to decrypt
	/*if user["password"] == vaultPwd {
		fmt.Println("Successfully logged")
	}*/

	//fmt.Println("Inserted a Single Record ", insertResult.InsertedID)
}
