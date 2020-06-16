package prometheus

import (
	"net/http"
	"time"

	"github.com/micro/go-micro/v2/client"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/server"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheus represents go micro monitoring plugin for prometheus
type Prometheus interface {
	NewClientWrapper() client.Wrapper
	NewCallWrapper() client.CallWrapper
	NewHandlerWrapper() server.HandlerWrapper
	NewSubscriberWrapper() server.SubscriberWrapper
}

type plugin struct {
	opts Options

	uptime *prometheus.CounterVec

	callCnt *prometheus.CounterVec
	callDur *prometheus.HistogramVec

	streamCnt *prometheus.CounterVec
	streamDur *prometheus.HistogramVec

	publishCnt *prometheus.CounterVec
	publishDur *prometheus.HistogramVec

	subscribeCnt *prometheus.CounterVec
	subscribeDur *prometheus.HistogramVec

	handleCnt *prometheus.CounterVec
	handleDur *prometheus.HistogramVec
}

// NewPrometheus generates a new set of metrics with a certain subsystem name
func NewPrometheus(opts ...Option) Prometheus {
	p := &plugin{opts: Options{ListenAddress: ":8080", Subsystem: "micro", MetricsPath: "/metrics"}}
	for _, o := range opts {
		o(&p.opts)
	}
	p.registerMetrics()
	go func() {
		for range time.Tick(time.Second) {
			p.uptime.WithLabelValues(p.opts.ServiceID, p.opts.ServiceName, p.opts.ServiceVersion).Inc()
		}
	}()
	go func() {
		http.Handle(p.opts.MetricsPath, promhttp.Handler())
		if err := http.ListenAndServe(p.opts.ListenAddress, nil); err != nil {
			log.Error(err)
		}
	}()
	return p
}

// Options represents prometheus metrics options
type Options struct {
	ServiceID      string
	ServiceName    string
	ServiceVersion string
	Subsystem      string
	ListenAddress  string
	MetricsPath    string
}

// Option represents prometheus options update method
type Option func(*Options)

// ServiceID represents service id
func ServiceID(id string) Option {
	return func(opts *Options) {
		opts.ServiceID = id
	}
}

// ServiceName represents service name
func ServiceName(name string) Option {
	return func(opts *Options) {
		opts.ServiceName = name
	}
}

// ServiceVersion represents service version
func ServiceVersion(version string) Option {
	return func(opts *Options) {
		opts.ServiceVersion = version
	}
}

// ListenAddress represents exposing metrics address.
func ListenAddress(address string) Option {
	return func(opts *Options) {
		opts.ListenAddress = address
	}
}

// MetricsPath represents exposing metrics path
func MetricsPath(path string) Option {
	return func(opts *Options) {
		opts.MetricsPath = path
	}
}
