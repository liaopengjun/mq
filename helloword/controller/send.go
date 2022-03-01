package controller

import (
	"github.com/gin-gonic/gin"
	"helloword/models/mq"
)

func Send(c *gin.Context)  {
	mq.SendHelloWorld()
}
