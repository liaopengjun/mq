package controller

import (
	"github.com/gin-gonic/gin"
	"work_queues/models/mq"
)

func Task(c *gin.Context) {
	s := "msg"
	for i := 0; i < 4; i++ {
		s = s + "."
		mq.SendHelloWorld(s)
	}
}
