package redis

import (
	"consumer-service/logs"
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Client *redis.Client
}

// NewRedis initializes a Redis client.
func NewRedis(addr string) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &Redis{Client: client}
}

// SaveUser saves a user record to Redis.
func (r *Redis) SaveUser(id int, user interface{}) error {
	// Convert user data to JSON
	userJSON, err := json.Marshal(user)
	if err != nil {
		logs.Log.Errorln("Failed to marshal user to JSON:", err)
		return err
	}

	// Store the JSON in Redis with the key "user:<id>"
	key := fmt.Sprintf("user:%d", id)
	err = r.Client.Set(context.Background(), key, userJSON, 0).Err()
	if err != nil {
		logs.Log.Errorf("Failed to save user in Redis with key %s: %v\n", key, err)
		return err
	}

	logs.Log.Infof("Successfully saved user in Redis with key: %s\n", key)
	return nil
}
