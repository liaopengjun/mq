package mq

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type Callback func(msg string)

func SendWork(body string) {
	mq,err := Connect()
	ch, err := mq.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ch.Close()
	// 3. 声明消息要发送到的队列
	q, err := ch.QueueDeclare(
		"Work", // 队列名称
		false,  // 定义是否持久化 true持久化队列 false不持久
		false,  // 是否自动删除
		false,  // 是否独占队列
		false,  // 是否不等待
		nil,    // 其他参数
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 4.将消息发布到声明的队列
	err = ch.Publish(
		"",     // 交换机
		q.Name, 		// 路由key
		false,  // 是否强制
		false,  // 即时
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, //持久
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf(" [x] Sent %s", body)
}

func ReceiveWork() {
	mq, err := Connect()
	if err != nil{
		fmt.Println(err)
		return
	}
	defer mq.Close()
	ch, err := mq.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"Work", // 队列名称
		false,  // 定义是否持久化 true持久化队列 false不持久
		false,  // 是否自动删除
		false,  // 是否独占队列
		false,  // 是否不等待
		nil,    // 其他参数
	)
	// 获取接收消息的Delivery通道
	msgs, err := ch.Consume(
		q.Name, 		 // 队列名称
		"",     // 消费者
		false,   // 自动确认
		false,  // 是否独占队列
		false,  // 非本地
		false,  // 是否不等待
		nil,    // 其他参数
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
			d.Ack(false) // 手动传递消息确认
		}
	}()
	<-forever
}

func SendCalculatework(body string)  {
	mq,err := Connect()
	ch, err := mq.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ch.Close()
	// 3. 声明消息要发送到的队列
	q, err := ch.QueueDeclare(
		"task_queue", // 队列名称
		true,  // 定义是否持久化 true持久化队列 false不持久
		false,  // 是否自动删除
		false,  // 是否独占队列
		false,  // 是否不等待
		nil,    // 其他参数
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 4. 然后我们可以将消息发布到声明的队列
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // 立即
		false,  // 强制
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // 持久
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	if err != nil {
		fmt.Printf("publish a message failed, err:%v\n", err)
		return
	}
	log.Printf(" [x] Sent %s", body)
}

func ReceiveCalculateWork()  {
	mq,err := Connect()
	if err != nil {
		fmt.Printf("connect to RabbitMQ failed, err:%v\n", err)
		return
	}
	defer mq.Close()

	ch, err := mq.Channel()
	if err != nil {
		fmt.Printf("open a channel failed, err:%v\n", err)
		return
	}
	defer ch.Close()

	// 声明一个queue
	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // 声明为持久队列
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	err = ch.Qos(
		1,     // 预计计数
		0,     // 预计大小
		false, // 是否全局
	)
	if err != nil {
		fmt.Printf("ch.Qos() failed, err:%v\n", err)
		return
	}

	// 立即返回一个Delivery的通道
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // 注意这里传false,关闭自动消息确认
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		fmt.Printf("ch.Consume failed, err:%v\n", err)
		return
	}

	// 开启循环不断地消费消息
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			d.Ack(false) // 手动传递消息确认
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
