package api

import (
	"consumer-service/connectors"

	"github.com/gorilla/mux"
)

func SetupRoutes(db *connectors.Postgres) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/users", GetUsersHandler(db)).Methods("GET")
	return router
}
