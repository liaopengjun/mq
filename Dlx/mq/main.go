package mq

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
)

type Callback func(msg string)

func Connect() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@49.234.9.118:5672/")
	return conn, err
}

func PublishDlx(exchangeA string, body string) error {
	//建立连接
	conn, err := Connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	//创建一个Channel
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	//消息发送到A交换机
	err = channel.Publish(exchangeA, "", false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         []byte(body),
	})
	fmt.Println("msg:",body)
	return err
}

func ConsumerDlx(exchangeA string, queueAName string, exchangeB string, queueBName string, ttl int, callback Callback) {
	//建立连接
	conn, err := Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	//创建一个Channel
	channel, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer channel.Close()

	//创建A交换机
	//创建A队列
	//A交换机和A队列绑定
	err = channel.ExchangeDeclare(
		exchangeA, // name
		"fanout",  // type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	//创建一个queue，指定消息过期时间，并且绑定过期以后发送到那个交换机
	queueA, err := channel.QueueDeclare(
		queueAName, // name
		true,       // durable
		false,      // delete when usused
		false,      // exclusive
		false,      // no-wait
		amqp.Table{
			// 当消息过期时把消息发送到 exchangeB
			"x-dead-letter-exchange": exchangeB,
			"x-message-ttl":          ttl,
			//"x-dead-letter-queue" : queueBName,
			//"x-dead-letter-routing-key" :
		},
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	//A交换机和A队列绑定
	err = channel.QueueBind(
		queueA.Name, // queue name
		"",          // routing key
		exchangeA,   // exchange
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	//创建B交换机
	//创建B队列
	//B交换机和B队列绑定
	err = channel.ExchangeDeclare(
		exchangeB, // name
		"fanout",  // type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	//创建一个queue
	queueB, err := channel.QueueDeclare(
		queueBName, // name
		true,       // durable
		false,      // delete when usused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	//B交换机和B队列绑定
	err = channel.QueueBind(
		queueB.Name, // queue name
		"",          // routing key
		exchangeB,   // exchange
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	msgs, err := channel.Consume(queueB.Name, "", false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			s := BytesToString(&(d.Body))
			callback(*s)
			d.Ack(false)
		}
	}()

	fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C\n")
	<-forever
}


func BytesToString(b *[]byte) *string {
	s := bytes.NewBuffer(*b)
	r := s.String()
	return &r
}