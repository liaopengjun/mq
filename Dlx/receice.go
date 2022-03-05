package main

import (
	"fmt"
	"lwzgo/mq/Dlx/mq"
)

func main()  {
	mq.ConsumerDlx("dlx_a","dlx_queus_a","dlx_b","dlx_queus_b",10000,callback)
}

func callback(s string) {
	fmt.Println(s)
}