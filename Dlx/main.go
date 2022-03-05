package main

import (
	"lwzgo/mq/Dlx/mq"
)

func main()  {
	mq.PublishDlx("dlx_a","hello")
}
