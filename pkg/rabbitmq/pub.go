package rabbitmq

import (
	"app-noti/common"
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type PublishClient interface {
	PublishMessage(ctx context.Context, channelPrefix string, exchangeName string, routingKey string, mandatory bool, immediate bool, message string) error
}

type publishClient struct {
	ctx    context.Context
	prefix string
	uri    string
	client *client
}

type client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queues  []queue
	channs  map[string]*channel
}

func NewPublishClient(prefix string) *publishClient {
	c := &client{channs: map[string]*channel{}}
	return &publishClient{prefix: prefix, client: c}
}

func (s *publishClient) Run() error {
	if err := s.configure(); err != nil {
		return err
	}

	for i, _ := range s.client.channs {
		chann, err := s.client.conn.Channel()
		if err != nil {
			return err
		}

		channCfg := s.client.channs[i]
		err = chann.ExchangeDeclare(channCfg.exchange.Name, channCfg.exchange.Kind, channCfg.exchange.Durable, channCfg.exchange.AutoDelete, channCfg.exchange.Internal, channCfg.exchange.NoWait, nil)
		if err != nil {
			log.Fatalf("exchange.declare: %v", err)
		}

		s.client.channs[channCfg.prefix].conn = chann
	}

	return nil
}

func (s *publishClient) InitChannels(ops ...*ChannOps) {
	for i, _ := range ops {
		op := ops[i]
		s.client.channs[op.Prefix] = &channel{prefix: op.Prefix, exchange: op.Exchange}
	}
}

func (s *publishClient) Get() interface{} {
	return s.client
}

func (s *publishClient) configure() error {
	s.initFlags()

	if s.prefix == "" {
		return errors.New(common.DataIsNullErr(s.prefix))
	}

	if s.uri == "" {
		return errors.New(common.DataIsNullErr(fmt.Sprintf("%v uri", s.GetPrefix())))
	}

	ctx := context.Background()

	s.ctx = ctx

	conn, err := amqp.Dial(s.uri)
	if err != nil {
		return err
	}
	s.client.conn = conn

	return nil
}

func (s *publishClient) initFlags() {
	uri := os.Getenv(common.ENV_RABBIT_URI)
	s.uri = uri
}

func (s *publishClient) GetPrefix() string {
	return s.prefix
}

func (s *publishClient) Stop() <-chan bool {
	stop := make(chan bool)
	go func() {
		s.client.conn.Close()
		stop <- true
	}()
	return stop
}

func (s *client) PublishMessage(ctx context.Context, channelPrefix string, exchangeName string, routingKey string, mandatory bool, immediate bool, message string) error {
	err := s.channs[channelPrefix].conn.PublishWithContext(ctx, exchangeName, routingKey,
		mandatory, // mandatory
		immediate, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

	return err
}
