package main

import (
	"github.com/gin-gonic/gin"
	"work_queues/controller"
)

func main() {
	r := gin.Default()
	r.GET("/new_task", controller.Task)
	r.GET("/new_task2", controller.CalculateWork)
	gin.SetMode(gin.ReleaseMode)
	r.Run(":9090")
}
