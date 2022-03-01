package controller

import (
	"github.com/gin-gonic/gin"
	"helloword/models/mq"
)

func Receive(c *gin.Context)  {
	mq.ReceiveHelloWorld()
}

//func callback(s string)  {
//	fmt.Printf("msg is: %s\n",s)
//}