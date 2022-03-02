package mq

import (
	"github.com/streadway/amqp"
)

func Connect() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@49.234.9.118:5672/")
	return conn, err
}

