package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

type Exchange struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
}

type channel struct {
	prefix   string
	exchange Exchange
	conn     *amqp.Channel
}

type ChannOps struct {
	Prefix   string
	Exchange Exchange
}
