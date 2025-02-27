package redis

import (
	"consumer-service/logs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveUser(t *testing.T) {
	logs.LoggerConfiguration()
	addr := "localhost:6379"
	r := NewRedis(addr)

	testUser := map[string]interface{}{
		"name": "John Doe",
		"age":  30,
	}

	err := r.SaveUser(1, testUser)
	assert.NoError(t, err)

	// Verify the user exists in Redis
	key := "user:1"
	val, err := r.Client.Get(r.Client.Context(), key).Result()
	assert.NoError(t, err)
	assert.NotEmpty(t, val)
}
