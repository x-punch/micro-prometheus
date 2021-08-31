package main

import (
	"github.com/asim/go-micro/v3"
	prometheus "github.com/x-punch/micro-prometheus/v3"
)

func main() {
	service := micro.NewService(micro.Name("go.micro.prom.testing"), micro.Version("0.0.0"))
	promOpts := []prometheus.Option{
		prometheus.ServiceID(service.Server().Options().Id),
		prometheus.ServiceName(service.Name()),
		prometheus.ServiceVersion(service.Server().Options().Version),
		prometheus.ListenAddress(":3000"),
	}
	prom := prometheus.NewPrometheus(promOpts...)
	microOpts := []micro.Option{
		micro.WrapHandler(prom.NewHandlerWrapper()),
		micro.WrapSubscriber(prom.NewSubscriberWrapper()),
	}
	service.Init(microOpts...)
	if err := service.Run(); err != nil {
		panic(err)
	}
}
