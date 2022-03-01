package mq

import (
	"github.com/streadway/amqp"
)

var mq *amqp.Connection
func Init()  (err error) {
	mq, err = amqp.Dial("amqp://guest:guest@49.234.9.118:5672/")
	return
}

//释放资源
func Close()  {
	_ = mq.Close()
}
