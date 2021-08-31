package prometheus

import (
	"github.com/asim/go-micro/v3/logger"
	"github.com/prometheus/client_golang/prometheus"
)

func (p *plugin) registerMetrics() {
	p.uptime = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "uptime",
		Subsystem: p.opts.Subsystem,
		Help:      "Service uptime",
	}, []string{"id", "name", "version"})

	p.callCnt = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "call_total",
		Subsystem: p.opts.Subsystem,
		Help:      "Call processed, partitioned by endpoint and status",
	}, []string{"id", "name", "version", "endpoint", "status"})
	p.callDur = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "call_duration_seconds",
		Subsystem: p.opts.Subsystem,
		Help:      "Call latencies in seconds, partitioned by endpoint",
	}, []string{"id", "name", "version", "endpoint"})

	p.streamCnt = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "stream_total",
		Subsystem: p.opts.Subsystem,
		Help:      "Stream processed, partitioned by endpoint and status",
	}, []string{"id", "name", "version", "endpoint", "status"})
	p.streamDur = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "stream_duration_seconds",
		Subsystem: p.opts.Subsystem,
		Help:      "Stream latencies in seconds, partitioned by endpoint",
	}, []string{"id", "name", "version", "endpoint"})

	p.publishCnt = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "publish_total",
		Subsystem: p.opts.Subsystem,
		Help:      "Publish processed, partitioned by endpoint and status",
	}, []string{"id", "name", "version", "endpoint", "status"})
	p.publishDur = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "publish_duration_seconds",
		Subsystem: p.opts.Subsystem,
		Help:      "Publish latencies in seconds, partitioned by endpoint",
	}, []string{"id", "name", "version", "endpoint"})

	p.subscribeCnt = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "subscribe_total",
		Subsystem: p.opts.Subsystem,
		Help:      "Subscribe processed, partitioned by endpoint and status",
	}, []string{"id", "name", "version", "endpoint", "status"})
	p.subscribeDur = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "subscribe_duration_seconds",
		Subsystem: p.opts.Subsystem,
		Help:      "Subscribe latencies in seconds, partitioned by endpoint",
	}, []string{"id", "name", "version", "endpoint"})

	p.handleCnt = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "handle_total",
		Subsystem: p.opts.Subsystem,
		Help:      "Handle processed, partitioned by endpoint and status",
	}, []string{"id", "name", "version", "endpoint", "status"})
	p.handleDur = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "handle_duration_seconds",
		Subsystem: p.opts.Subsystem,
		Help:      "Handle latencies in seconds, partitioned by endpoint",
	}, []string{"id", "name", "version", "endpoint"})

	for _, collector := range []prometheus.Collector{p.uptime, p.callCnt, p.callDur, p.streamCnt, p.streamDur, p.publishCnt, p.publishDur, p.handleCnt, p.handleDur, p.subscribeCnt, p.subscribeDur} {
		if err := prometheus.DefaultRegisterer.Register(collector); err != nil {
			logger.Error(err)
		}
	}
}
