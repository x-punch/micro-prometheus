# Micro Prometheus
Go Micro monitoring plugin, export prometheus metrics through target address.

## Usage
```go
import (
	"github.com/micro/go-micro/v2"
	prometheus "github.com/x-punch/micro-prometheus/v2"
)

func main(){
    service := micro.NewService()
    prom := prometheus.NewPrometheus()
    service.Init(micro.WrapHandler(prom.NewHandlerWrapper()), micro.WrapSubscriber(prom.NewSubscriberWrapper()))
}
```
```shell
curl http://:8080/metrics
```

## Tips
```
When plugin created, it'll listen on http address(default :8080), if the port is not avaiable or other issues, you'll have an error, but the program'll still be running.
```