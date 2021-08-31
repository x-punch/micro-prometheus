# Micro Prometheus
Go Micro monitoring plugin, export prometheus metrics through target address.

## Usage
```go
import (
	"github.com/asim/go-micro/v3"
	prometheus "github.com/x-punch/micro-prometheus/v3"
)

func main() {
	service := micro.NewService(micro.Name("go.micro.prometheus.testing"), micro.Version("1.0.0"))
	promOpts := []prometheus.Option{
		prometheus.ServiceID(service.Server().Options().Id),
		prometheus.ServiceName(service.Name()),
		prometheus.ServiceVersion(service.Server().Options().Version),
		prometheus.ListenAddress(":8080"),
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
```
```shell
curl http://:8080/metrics
```

## Tips
```
When plugin created, it'll listen on http address(default :8080), if the port is not avaiable or other issues, you'll have an error, but the program'll still be running.
```