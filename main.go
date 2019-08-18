package main

import (
	"github.com/SuperCLine/wxpay/wxpay"
)

func main()  {

	srv := wxpay.NewDefaultPayService()
	srv.Start()
}
