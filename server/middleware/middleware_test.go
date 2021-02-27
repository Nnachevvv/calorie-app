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
		Password: "test-password",
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
	mockMongoDB.EXPECT().Find("test-name", gomock.Any()).Return(bson.M{}, nil).Times(1)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(middleware.LoginToSystem)
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	// Check the response body is what we expect.
	expected := `{"_id":"000000000000000000000000","username":"test-name","password":"test-password"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
