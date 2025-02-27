package main

import (
	"net/http"
	"os"
	"producer-service/api"
	"producer-service/logs"
	"producer-service/rabbitmq"
	"producer-service/helpers"

	"github.com/spf13/viper"
)

func init() {
	//Initialize the application configurations
	applicationPath, _ := os.Getwd()
	helpers.Configuration(applicationPath)

	//Initialize the log configurations
	logs.LoggerConfiguration()
}

func main() {
	// Initialize RabbitMQ producer
	producer, err := rabbitmq.NewProducer(viper.GetString("rabbitmq_url"))
	if err != nil {
		logs.Log.Fatalf("Failed to connect to RabbitMQ: %v\n", err)
	}
	defer producer.Close()

	// Set up routes
	router := api.SetupRoutes(producer)

	// Start the server
	logs.Log.Infoln("producer app is running on http://localhost:8080")
	logs.Log.Fatalln(http.ListenAndServe(":8080", router))
}
