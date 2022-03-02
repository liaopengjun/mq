package controller

import (
	"github.com/gin-gonic/gin"
	"work_queues/models/mq"
)

func Task(c *gin.Context) {
	s := "msg"
	for i := 0; i < 1; i++ {
		s = s + "........."
		mq.SendWork(s)
	}
}

func CalculateWork(c *gin.Context) {
	s := "msg"
	for i := 0; i < 10; i++ {
		s = s + "."
		mq.SendCalculatework(s)
	}
}
