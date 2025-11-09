package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitProducer struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	logger   *logrus.Logger
}

// NewRabbitProducer creates a new producer connected to RabbitMQ.
func NewRabbitProducer(amqpURL, exchange string, logger *logrus.Logger) (*RabbitProducer, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	err = ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	logger.Infof("✅ RabbitMQ producer connected (exchange=%s)", exchange)

	return &RabbitProducer{
		conn:     conn,
		channel:  ch,
		exchange: exchange,
		logger:   logger,
	}, nil
}

// Produce publishes a JSON-encoded message to the exchange using eventType as routing key.
func (p *RabbitProducer) Produce(ctx context.Context, eventType string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		p.logger.Errorf("failed to marshal message: %v", err)
		return err
	}

	err = p.channel.PublishWithContext(
		ctx,
		p.exchange,
		eventType,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		p.logger.Errorf("failed to publish message: %v", err)
		return err
	}

	p.logger.Infof("[Producer] event=%s message=%s", eventType, string(body))
	return nil
}

func (p *RabbitProducer) Close() {
	_ = p.channel.Close()
	_ = p.conn.Close()
	p.logger.Info("RabbitMQ producer closed")
}
