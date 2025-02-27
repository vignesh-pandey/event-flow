package api

import (
	"producer-service/rabbitmq"

	"github.com/gorilla/mux"
)

func SetupRoutes(producer *rabbitmq.Producer) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/upload-csv", UploadCSV(producer)).Methods("POST")
	return router
}
