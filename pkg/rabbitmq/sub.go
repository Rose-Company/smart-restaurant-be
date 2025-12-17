package rabbitmq

import (
	"app-noti/common"
	"context"
	"errors"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type publisher struct {
	exchange string
	queue    string
	routing  string
}

type rabbitmq struct {
	prefix     string
	conn       *amqp.Connection
	uri        string
	consumers  []consumer
	publishers map[string]*publisher
}

func NewRabbitMQ(prefix string, uri string) *rabbitmq {
	return &rabbitmq{prefix: prefix, uri: uri, publishers: map[string]*publisher{}}
}

func (s *rabbitmq) Run() error {
	if err := s.configure(); err != nil {
		return err
	}

	var forever chan struct{}

	go func() {
		for i, _ := range s.consumers {
			fmt.Println(fmt.Sprintf("Consumer %v is running", s.consumers[i].GetPrefix(context.Background())))

			go func(cons *consumer, conn *amqp.Connection) {
				ctx := context.Background()
				defer fmt.Println(fmt.Sprintf("Consumer %v is stopped", cons.GetPrefix(context.Background())))
				//defer conn.Close()
				chann, err := conn.Channel()
				if err != nil {
					log.Panic(err)
					return
				}

				if err := cons.Run(ctx, chann); err != nil {
					log.Panic(err)
					return
				}
			}(&s.consumers[i], s.conn)
		}

		<-forever
	}()

	return nil
}

func (s *rabbitmq) InitConsumers(ctx context.Context, consumers ...any) {
	for i, _ := range consumers {
		s.consumers = append(s.consumers, consumers[i].(consumer))
	}
}

func (s *rabbitmq) GetConn() *amqp.Connection {
	return s.conn
}

func (s *rabbitmq) configure() error {
	if s.prefix == "" {
		return errors.New(common.DataIsNullErr(s.prefix))
	}

	if s.uri == "" {
		return errors.New(common.DataIsNullErr(fmt.Sprintf("%v uri", s.GetPrefix())))
	}

	conn, err := amqp.Dial(s.uri)
	if err != nil {
		return err
	}
	s.conn = conn

	return nil
}

func (s *rabbitmq) GetPrefix() string {
	return s.prefix
}

func (s *rabbitmq) Stop() <-chan bool {
	//s.conn.Close()
	stop := make(chan bool)
	go func() {
		stop <- true
	}()
	return stop
}

func (s *rabbitmq) InitPublisher(key string, exchange string, routing string) {
	ps := s.publishers

	pubConfig := publisher{
		exchange: exchange,
		queue:    routing,
	}
	ps[key] = &pubConfig
	s.publishers = ps

}

func (s *rabbitmq) Get() interface{} {
	return s
}
