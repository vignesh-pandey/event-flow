package rabbitmq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProducer(t *testing.T) {
	// Replace with a local RabbitMQ server or mock
	url := "amqp://guest:guest@localhost:5672/"
	producer, err := NewProducer(url)
	assert.NoError(t, err)
	assert.NotNil(t, producer)

	// Close the producer after test
	producer.Close()
}

func TestPublish(t *testing.T) {
	url := "amqp://guest:guest@localhost:5672/"
	producer, err := NewProducer(url)
	assert.NoError(t, err)
	defer producer.Close()

	err = producer.Publish("test_queue", "test_message")
	assert.NoError(t, err)
}
