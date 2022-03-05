package main

import "PublishSubscribe/mq"

func main()  {
	sendMsg :="log"
	mq.Publish(sendMsg)
}
