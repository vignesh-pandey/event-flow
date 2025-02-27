package rabbitmq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConsumer(t *testing.T) {
	url := "amqp://guest:guest@localhost:5672/"
	consumer, err := NewConsumer(url, "test_queue")
	assert.NoError(t, err)
	assert.NotNil(t, consumer)

	// Close the consumer after test
	consumer.Close()
}

