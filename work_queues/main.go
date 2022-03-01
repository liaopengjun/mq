package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"work_queues/controller"
	"work_queues/models/mq"
)

func main() {
	r := gin.Default()
	if err := mq.Init(); err != nil {
		log.Fatalf("connect mq failed err :%s", err)
	}
	defer mq.Close()
	r.GET("/new_task", controller.Task)
	gin.SetMode(gin.ReleaseMode)
	r.Run(":9090")
}
