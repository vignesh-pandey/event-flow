package rabbitmq

import (
	"producer-service/logs"

	"github.com/streadway/amqp"
)

type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewProducer(url string) (*Producer, error) {
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

	return &Producer{conn: conn, channel: ch}, nil
}

func (p *Producer) Publish(queueName string, message string) error {
	_, err := p.channel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		logs.Log.Errorln("Failed to declare a queue:", err)
		return err
	}

	return p.channel.Publish("", queueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	})
}

func (p *Producer) Close() {
	if err := p.channel.Close(); err != nil {
		logs.Log.Infoln("Failed to close RabbitMQ channel:", err)
	}
	if err := p.conn.Close(); err != nil {
		logs.Log.Infoln("Failed to close RabbitMQ connection:", err)
	}
}
