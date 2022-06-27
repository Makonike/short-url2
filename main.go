package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"short-url/controller"
	"short-url/object"
)

func main() {
	h := server.Default()
	controller.ShortUrl(h)
	h.Spin()
}

func init() {
	hl := hlog.DefaultLogger()
	err := object.SetupSetting()
	object.InitAdapter()
	if err != nil {
		hl.Error("Init Server Error")
	}
}
