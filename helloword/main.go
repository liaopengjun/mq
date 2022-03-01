package main

import (
	"github.com/gin-gonic/gin"
	"helloword/controller"
	"helloword/models/mq"
	"log"
)

func main(){
	r := gin.Default()
	if err := mq.Init();err !=nil{
		log.Fatalf("connect mq failed err :%s", err)
	}
	defer mq.Close()
	r.GET("/send", controller.Send)
	r.GET("/receive", controller.Receive)
	gin.SetMode(gin.ReleaseMode)
	r.Run(":9090")
}
