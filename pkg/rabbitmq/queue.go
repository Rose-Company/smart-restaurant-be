package rabbitmq

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
)

type queue struct {
	prefix     string
	exchange   string
	queueName  string
	routingKey string
	instance   *amqp091.Queue
}

func NewQueue(prefix string, exchange string, queueName string, routingKey string) *queue {
	return &queue{
		prefix:     prefix,
		exchange:   exchange,
		queueName:  queueName,
		routingKey: routingKey,
	}
}

func (c *queue) GetPrefix(ctx context.Context) string {
	return c.prefix
}

func (c *queue) init(ctx context.Context, chann *amqp091.Channel) error {
	q, err := chann.QueueDeclare(c.queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	c.instance = &q

	return nil
}
