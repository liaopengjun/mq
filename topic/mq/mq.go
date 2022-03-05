package mq


import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strings"
	"time"
)

func Connect() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@49.234.9.118:5672/")
	return conn, err
}

func SendWork()  {
	mq,err := Connect()
	ch,err := mq.Channel()
	if err != nil{
		fmt.Println(err)
	}
	defer ch.Close()
	// 3. 声明消息要发送到的队列
	err = ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	body := bodyFrom(os.Args)
	// 4.将消息发布到声明的队列
	err = ch.Publish (
		"logs_topic",         // exchange
		severityFrom(os.Args), 			// routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(body),
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
	err = ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [info] [warning] [error]", os.Args[0])
		os.Exit(0)
	}
	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, "logs_direct", s)
		err = ch.QueueBind(
			q.Name,        // queue name
			s,             // routing key
			"logs_topic", // exchange
			false,
			nil)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
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

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 3) || os.Args[2] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[2:], " ")
	}
	return s
}

func severityFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "anonymous.info"
	} else {
		s = os.Args[1]
	}
	return s
}