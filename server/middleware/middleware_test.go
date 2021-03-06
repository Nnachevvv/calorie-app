package middleware_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nnachevv/calorieapp/models"
	"github.com/Nnachevv/calorieapp/server/middleware"
	"github.com/Nnachevv/calorieapp/server/middleware/mocks"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson"
)

func TestLoginToSystem(t *testing.T) {
	loginRequest := models.User{
		Username: "test-name",
		Password: "ValidPassword1@",
	}

	requestByte, _ := json.Marshal(loginRequest)
	requestReader := bytes.NewReader(requestByte)

	req, err := http.NewRequest("POST", "/login", requestReader)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMongoDB := mocks.NewMockMongoDatabase(mockCtrl)
	middleware.MongoService = mockMongoDB
	mockMongoDB.EXPECT().Find("test-name").Return(bson.M{"password": "ValidPassword1@"}, nil).Times(1)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(middleware.LoginUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

}

func TestUserIsNotFoundInSystem(t *testing.T) {
	loginRequest := models.User{
		Username: "test-name",
		Password: "ValidPassword1@",
	}

	requestByte, _ := json.Marshal(loginRequest)
	requestReader := bytes.NewReader(requestByte)

	req, err := http.NewRequest("POST", "/login", requestReader)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMongoDB := mocks.NewMockMongoDatabase(mockCtrl)
	middleware.MongoService = mockMongoDB

	mockMongoDB.EXPECT().Find("test-name").Return(nil, middleware.ErrUserIsNotFound).Times(1)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(middleware.LoginUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	if rr.Body.String() != middleware.ErrUserIsNotFound.Error() {
		t.Errorf("handler returned wrong error: got %v want %w", rr.Body.String(), middleware.ErrUserIsNotFound)
	}
}

func TestUserTypeWrongHisPassword(t *testing.T) {
	loginRequest := models.User{
		Username: "test-name",
		Password: "test-wrong-password",
	}

	requestByte, _ := json.Marshal(loginRequest)
	requestReader := bytes.NewReader(requestByte)

	req, err := http.NewRequest("POST", "/login", requestReader)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMongoDB := mocks.NewMockMongoDatabase(mockCtrl)
	middleware.MongoService = mockMongoDB

	mockMongoDB.EXPECT().Find("test-name").Return(bson.M{"password": "ValidPassword1@"}, nil).Times(1)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(middleware.LoginUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	if rr.Body.String() != middleware.ErrWrongPassword.Error() {
		t.Errorf("handler returned wrong error: got %v want %w", rr.Body.String(), middleware.ErrWrongPassword)
	}
}

func TestUserRegister(t *testing.T) {
	loginRequest := models.RegisterUser{
		Username:        "test-name",
		Password:        "ValidPassword1@",
		ConfirmPassword: "ValidPassword1@",
		Email:           "test@mail.com",
	}

	requestByte, _ := json.Marshal(loginRequest)
	requestReader := bytes.NewReader(requestByte)

	req, err := http.NewRequest("POST", "/register", requestReader)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMongoDB := mocks.NewMockMongoDatabase(mockCtrl)
	middleware.MongoService = mockMongoDB

	mockMongoDB.EXPECT().Find("test-name").Return(nil, middleware.ErrUserIsNotFound).Times(1)
	mockMongoDB.EXPECT().Add(loginRequest).Return(nil).Times(1)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(middleware.RegisterUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
func TestAlreadyExistingUser(t *testing.T) {
	loginRequest := models.RegisterUser{
		Username:        "test-name",
		Password:        "ValidPassword1@",
		ConfirmPassword: "ValidPassword1@",
		Email:           "test@mail.com",
	}

	requestByte, _ := json.Marshal(loginRequest)
	requestReader := bytes.NewReader(requestByte)

	req, err := http.NewRequest("POST", "/register", requestReader)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMongoDB := mocks.NewMockMongoDatabase(mockCtrl)
	middleware.MongoService = mockMongoDB

	mockMongoDB.EXPECT().Find("test-name").Return(bson.M{"user": "exist"}, nil).Times(1)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(middleware.RegisterUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusConflict)
	}

	if rr.Body.String() != middleware.ErrUserAlreadyExist.Error() {
		t.Errorf("handler returned wrong error: got %v want %w", rr.Body.String(), middleware.ErrUserAlreadyExist)
	}
}

func TestWhenPassedInvalidUserErrorIsThrown(t *testing.T) {
	invalidRegisterUsers := []models.RegisterUser{{
		Username:        "small",
		Password:        "ValidPassword1@",
		ConfirmPassword: "ValidPassword1@",
		Email:           "test@mail.com",
	}, {
		Username:        "12withnumbers",
		Password:        "ValidPassword1@",
		ConfirmPassword: "ValidPassword1@",
		Email:           "test@mail.com",
	}, {
		Username:        "_underscore",
		Password:        "ValidPassword1@",
		ConfirmPassword: "ValidPassword1@",
		Email:           "test@mail.com",
	},
		{
			Username:        "astringabove20characterslength",
			Password:        "ValidPassword1@",
			ConfirmPassword: "ValidPassword1@",
			Email:           "test@mail.com",
		},
		{
			Username:        ".dotstring",
			Password:        "ValidPassword1@",
			ConfirmPassword: "ValidPassword1@",
			Email:           "test@mail.com",
		},
		{
			Username:        "dots...attthemiddle",
			Password:        "ValidPassword1@",
			ConfirmPassword: "ValidPassword1@",
			Email:           "test@mail.com",
		},
		{
			Username:        "underscore__atmidle",
			Password:        "ValidPassword1@",
			ConfirmPassword: "ValidPassword1@",
			Email:           "test@mail.com",
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMongoDB := mocks.NewMockMongoDatabase(mockCtrl)
	middleware.MongoService = mockMongoDB

	for _, user := range invalidRegisterUsers {
		requestByte, _ := json.Marshal(user)
		requestReader := bytes.NewReader(requestByte)

		req, err := http.NewRequest("POST", "/register", requestReader)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(middleware.RegisterUser)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusUnprocessableEntity {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnprocessableEntity)
		}

		if rr.Body.String() != middleware.ErrUsernameIsInvalid.Error() {
			t.Errorf("handler returned wrong error: got %v want %w", rr.Body.String(), middleware.ErrUserAlreadyExist)
		}
	}

}

func TestPasswordIsInvalid(t *testing.T) {
	invalidUsersPasswords := []models.RegisterUser{{
		Username:        "test-name",
		Password:        "ShortWithoutNumber@",
		ConfirmPassword: "ShortWithoutNumber",
		Email:           "test@mail.com",
	}, {
		Username:        "test-name",
		Password:        "short@",
		ConfirmPassword: "short",
		Email:           "test@mail.com",
	}, {
		Username:        "test-name",
		Password:        "WithoutLatter@",
		ConfirmPassword: "WithoutLatter",
		Email:           "test@mail.com",
	},
		{
			Username:        "test-name",
			Password:        "notsamepassword",
			ConfirmPassword: "notsamepassword",
			Email:           "test@mail.com",
		},
		{
			Username:        "test-name",
			Password:        "withoutSpecialNumb",
			ConfirmPassword: "withoutSpecialNumb",
			Email:           "test@mail.com",
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMongoDB := mocks.NewMockMongoDatabase(mockCtrl)
	middleware.MongoService = mockMongoDB

	for _, user := range invalidUsersPasswords {
		requestByte, _ := json.Marshal(user)
		requestReader := bytes.NewReader(requestByte)

		req, err := http.NewRequest("POST", "/register", requestReader)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(middleware.RegisterUser)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusUnprocessableEntity {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnprocessableEntity)
		}

		if rr.Body.String() != middleware.ErrUserPasswordIsInvalid.Error() {
			t.Errorf("handler returned wrong error: got %v want %w", rr.Body.String(), middleware.ErrUserPasswordIsInvalid)
		}
	}
}
