package user_test

import (
	"testing"

	"github.com/Nnachevv/calorieapp/user"
)

func TestNewWithValid(t *testing.T) {
	got := user.New("username", "Password1")
	expected := user.User{"username", "Password1"}
	if got != expected {
		t.Errorf("New(username,password) ; want username password; got %s", got)
	}
}
