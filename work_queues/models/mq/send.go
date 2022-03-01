package mq

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type Callback func(msg string)

func SendHelloWorld(body string) {
	ch, err := mq.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ch.Close()
	// 3. 声明消息要发送到的队列
	q, err := ch.QueueDeclare(
		"Work", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 4.将消息发布到声明的队列
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf(" [x] Sent %s", body)
}

func ReceiveHelloWorld() {
	mq, err := amqp.Dial("amqp://guest:guest@49.234.9.118:5672/")
	ch, err := mq.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"Work", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	// 获取接收消息的Delivery通道
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dot_count := bytes.Count(d.Body, []byte(".")) // 数一下有几个.
			log.Printf("count: %d", dot_count)
			t := time.Duration(dot_count)
			time.Sleep(t * time.Second) // 模拟耗时的任务
			log.Printf("Done")
		}
	}()
	<-forever
}
