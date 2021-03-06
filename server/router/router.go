package router

import (
	"github.com/Nnachevv/calorieapp/server/middleware"
	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/login", middleware.LoginUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/register", middleware.RegisterUser).Methods("POST", "OPTIONS")

	return router
}
