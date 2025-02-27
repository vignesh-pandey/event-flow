package main

import (
	"consumer-service/api"
	"consumer-service/connectors"
	"consumer-service/logs"
	"consumer-service/helpers"
	"net/http"
	"os"
)

func init() {

	//Initialize the log configurations
	logs.LoggerConfiguration()

	//Initialize the application configurations
	applicationPath, _ := os.Getwd()
	helpers.Configuration(applicationPath)

	// Initialize the database connection instance
	connectors.NewPostgres()

	// Initialize the consumer
	go api.ConsumeData()
}

func main() {
	// Initialize PostgreSQL connection (if required for other routes)
	postgres := connectors.GetDatabaseConnection()

	// Set up routes
	router := api.SetupRoutes(postgres)

	// Start the server
	logs.Log.Infoln("Consumer app is running on http://localhost:8081")
	logs.Log.Fatalln(http.ListenAndServe(":8081", router))
}
