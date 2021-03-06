package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"unicode"

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

//ErrUserAlreadyExist represents when given user already exist
var ErrUserAlreadyExist = errors.New("user already exist")

//ErrUserPasswordIsInvalid represents when password is invalid
var ErrUserPasswordIsInvalid = errors.New("invalid password")

//ErrUsernameIsInvalid represent when username is invalid
var ErrUsernameIsInvalid = errors.New("invalid password")

// LoginUser logins user into the system
func LoginUser(w http.ResponseWriter, r *http.Request) {
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

// RegisterUser register user into the system
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var registerUserRegquest models.RegisterUser
	err := json.NewDecoder(r.Body).Decode(&registerUserRegquest)
	if err != nil {
		log.Fatal(err)
	}

	err = registerUser(registerUserRegquest)
	if err == ErrUserAlreadyExist {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(ErrUserAlreadyExist.Error()))
	}

	if err == ErrUsernameIsInvalid {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(ErrUsernameIsInvalid.Error()))
	}

	if err == ErrUserPasswordIsInvalid {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(ErrUserPasswordIsInvalid.Error()))
	}

}

// Check if user is in database
func verifyUser(login models.User) error {
	user, err := MongoService.Find(login.Username)
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

func verifyPassword(s string) bool {
	letters := 0
	var number, upper, special bool
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			letters++
		default:
			return false
		}
	}

	fmt.Println(number, upper, special, letters)
	return (number && upper && special && letters >= 7)

}

func registerUser(login models.RegisterUser) error {
	usernameRegex := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$`)
	if !usernameRegex.MatchString(login.Username) || len(login.Username) < 8 || len(login.Username) > 22 {
		return ErrUsernameIsInvalid
	}

	if login.Password != login.ConfirmPassword || !verifyPassword(login.Password) {
		return ErrUserPasswordIsInvalid
	}

	if _, err := MongoService.Find(login.Username); err == nil {
		return ErrUserAlreadyExist
	}

	if err := MongoService.Add(login); err != nil {
		return err
	}

	return nil
}
