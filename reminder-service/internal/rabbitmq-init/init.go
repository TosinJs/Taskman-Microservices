package rabbitmqinit

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func Init(URI string) (*amqp.Connection, error) {
	return amqp.Dial(URI)
}
