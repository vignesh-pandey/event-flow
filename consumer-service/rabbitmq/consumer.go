package rabbitmq

import (
	"consumer-service/logs"

	"github.com/streadway/amqp"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewConsumer(url, queueName string) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		logs.Log.Errorln("Failed to connect to RabbitMQ:", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		logs.Log.Errorln("Failed to open a channel:", err)
		return nil, err
	}

	_, err = ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		logs.Log.Errorln("Failed to declare a queue:", err)
		return nil, err
	}

	return &Consumer{conn: conn, channel: ch}, nil
}

func (c *Consumer) Consume(queueName string, handleMessage func(msg string)) {
	msgs, err := c.channel.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		logs.Log.Fatalf("Failed to consume messages: %v\n", err)
	}

	for msg := range msgs {
		handleMessage(string(msg.Body))
	}
}

func (c *Consumer) Close() {
	if err := c.channel.Close(); err != nil {
		logs.Log.Infoln("Failed to close RabbitMQ channel:", err)
	}
	if err := c.conn.Close(); err != nil {
		logs.Log.Infoln("Failed to close RabbitMQ connection:", err)
	}
}
