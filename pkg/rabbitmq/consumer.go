package rabbitmq

import (
	"app-noti/server"
	"context"
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type handler func(srv server.ServerContext, msg []byte) error

type consumer struct {
	prefix       string
	exchange     string
	exchangeType string
	queueName    string
	routingKey   string
	hdl          handler
	svr          server.ServerContext
	autoAck      bool
}

func NewConsumer(prefix string, exchange string, exchangeType string, queueName string, routingKey string, hdl handler, svr server.ServerContext, autoAck bool) *consumer {
	return &consumer{
		prefix:       prefix,
		exchange:     exchange,
		exchangeType: exchangeType,
		queueName:    queueName,
		routingKey:   routingKey,
		hdl:          hdl,
		svr:          svr,
		autoAck:      autoAck,
	}
}

func (c *consumer) GetPrefix(ctx context.Context) string {
	return c.prefix
}

func (c *consumer) GetDeathCount(table amqp091.Table) int64 {
	if h, ok := table["x-death"]; ok {
		hin := h.([]interface{})
		if len(hin) > 0 {
			i := hin[0].(amqp091.Table)
			if c, ok := i["count"]; ok {
				return c.(int64)
			}
		}
	}

	//count := h[0].(amqp091.Table)["count"]
	//count := h[0]["count"].(int)
	return 0
}

func (c *consumer) Run(ctx context.Context, chann *amqp091.Channel) error {

	dlxExchangeName := fmt.Sprintf("%v-%v-dlx", c.exchange, c.queueName)
	dlxQueueName := fmt.Sprintf("%v-queue-dlx", c.queueName)

	err := chann.ExchangeDeclare(
		c.exchange,     // name
		c.exchangeType, // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	q, err := chann.QueueDeclare(c.queueName, false, false, false, false, amqp091.Table{"x-dead-letter-exchange": dlxExchangeName})
	if err != nil {
		return err
	}

	err = chann.QueueBind(q.Name, c.routingKey, c.exchange, false, nil)
	if err != nil {
		return err
	}

	// dead letter exchange
	if !c.autoAck {
		err = chann.ExchangeDeclare(dlxExchangeName, "fanout", true, false, false, false, nil)
		_, err = chann.QueueDeclare(dlxQueueName, false, false, false, false, amqp091.Table{"x-dead-letter-exchange": c.exchange, "x-message-ttl": 10000})
		err = chann.QueueBind(dlxQueueName, c.routingKey, dlxExchangeName, false, nil)
	}

	var forever chan struct{}

	msgs, err := chann.Consume(
		q.Name,    // queue
		"",        // consumer
		c.autoAck, // auto ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       // args
	)

	go func() {
		for d := range msgs {
			err := c.hdl(c.svr, d.Body)
			if c.autoAck {
				if err != nil {
					fmt.Println(err)
				}
			} else {
				if err != nil {
					if c.GetDeathCount(d.Headers) > 3 {
						fmt.Println("expired after 3 retry", string(d.Body))
						d.Ack(false)
					}
					err := d.Nack(false, false)
					if err != nil {
						log.Println("Error: ", err.Error())
						return
					}

				} else {
					err := d.Ack(false)
					if err != nil {
						log.Println("Error: ", err.Error())
						return
					}
				}
			}

		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever

	return nil
}
