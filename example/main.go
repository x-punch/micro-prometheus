package main

import (
	"github.com/micro/go-micro/v2"
	prometheus "github.com/x-punch/micro-prometheus/v2"
)

func main() {
	service := micro.NewService(micro.Name("go.micro.prom.testing"), micro.Version("0.0.0"))
	prom := prometheus.NewPrometheus(prometheus.ServiceID(service.Server().Options().Id), prometheus.ServiceName(service.Name()), prometheus.ServiceVersion(service.Server().Options().Version), prometheus.ListenAddress(":3000"))
	service.Init(micro.WrapHandler(prom.NewHandlerWrapper()), micro.WrapSubscriber(prom.NewSubscriberWrapper()))
	if err := service.Run(); err != nil {
		panic(err)
	}
}
