package main

import "work_queues/models/mq"

func main() {
	mq.ReceiveCalculateWork()
}
