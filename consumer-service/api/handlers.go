package api

import (
	"consumer-service/connectors"
	"consumer-service/helpers"
	"consumer-service/logs"
	"consumer-service/rabbitmq"
	"consumer-service/redis"
	"encoding/json"
	"net/http"

	"github.com/spf13/viper"
)

func GetUsersHandler(db *connectors.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse filters from the query parameters
		filters := map[string][]string{}
		for key, values := range r.URL.Query() {
			filters[key] = values
		}

		// Fetch filtered users from the database
		users, err := db.GetFilteredUsers(filters)
		if err != nil {
			logs.Log.Errorln("Failed to fetch users:", err)
			http.Error(w, "Failed to fetch users: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with the users in JSON format
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(users); err != nil {
			logs.Log.Errorln("Failed to encode users to JSON:", err)
			http.Error(w, "Failed to encode users to JSON: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func ConsumeData() {
	consumer, err := rabbitmq.NewConsumer(viper.GetString("rabbitmq_url"), "csv_queue")
	if err != nil {
		logs.Log.Fatalf("Failed to connect to RabbitMQ: %v\n", err)
	}
	defer consumer.Close()

	dbConn := connectors.GetDatabaseConnection()
	redisClient := redis.NewRedis(viper.GetString("redis_address"))

	consumer.Consume("csv_queue", func(msg string) {
		// Decrypt the message
		encryptionKey := viper.GetString("encryption_key")
		decryptedData, err := helpers.Decrypt(msg, encryptionKey)
		if err != nil {
			logs.Log.Errorf("Failed to decrypt message: %v\n", err)
			return
		}
		// Parse the message into the User struct
		user := connectors.User{}
		if err := json.Unmarshal([]byte(decryptedData), &user); err != nil {
			logs.Log.Errorf("Failed to parse message: %v\n", err)
			return
		}

		// Insert the user data into the database
		err = dbConn.InsertUser(&user)
		if err != nil {
			logs.Log.Errorf("Failed to insert user into database: %v\n", err)
			return
		}

		// Save the user data into Redis
		err = redisClient.SaveUser(user.ID, user)
		if err != nil {
			logs.Log.Errorf("Failed to save user in Redis: %v\n", err)
			return
		}

		logs.Log.Infof("Successfully inserted user: %+v\n", user)

	})
}
