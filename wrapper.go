package prometheus

import (
	"context"
	"fmt"

	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
	"github.com/prometheus/client_golang/prometheus"
)

type wrapper struct {
	client.Client
	prom     *plugin
	callFunc client.CallFunc
}

func (p *plugin) NewClientWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		handler := &wrapper{Client: c, prom: p}
		return handler
	}
}

func (p *plugin) NewCallWrapper() client.CallWrapper {
	return func(fn client.CallFunc) client.CallFunc {
		handler := &wrapper{callFunc: fn, prom: p}
		return handler.CallFunc
	}
}

func (p *plugin) NewHandlerWrapper() server.HandlerWrapper {
	handler := &wrapper{prom: p}
	return handler.HandlerFunc
}

func (p *plugin) NewSubscriberWrapper() server.SubscriberWrapper {
	handler := &wrapper{prom: p}
	return handler.SubscriberFunc
}

func (w *wrapper) CallFunc(ctx context.Context, node *registry.Node, req client.Request, rsp interface{}, opts client.CallOptions) error {
	status, endpoint := "success", fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())

	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(s float64) {
		w.prom.callDur.WithLabelValues(w.prom.opts.ServiceID, w.prom.opts.ServiceName, w.prom.opts.ServiceVersion, endpoint).Observe(s)
	}))
	defer timer.ObserveDuration()

	err := w.callFunc(ctx, node, req, rsp, opts)
	if err != nil {
		status = "failure"
	}
	w.prom.callCnt.WithLabelValues(w.prom.opts.ServiceID, w.prom.opts.ServiceName, w.prom.opts.ServiceVersion, endpoint, status).Inc()

	return err
}

func (w *wrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	status, endpoint := "success", fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())

	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(s float64) {
		w.prom.callDur.WithLabelValues(w.prom.opts.ServiceID, w.prom.opts.ServiceName, w.prom.opts.ServiceVersion, endpoint).Observe(s)
	}))
	defer timer.ObserveDuration()

	err := w.Client.Call(ctx, req, rsp, opts...)
	if err != nil {
		status = "failure"
	}
	w.prom.callCnt.WithLabelValues(w.prom.opts.ServiceID, w.prom.opts.ServiceName, w.prom.opts.ServiceVersion, endpoint, status).Inc()

	return err
}

func (w *wrapper) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	status, endpoint := "success", fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())

	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(s float64) {
		w.prom.streamDur.WithLabelValues(w.prom.opts.ServiceID, w.prom.opts.ServiceName, w.prom.opts.ServiceVersion, endpoint).Observe(s)
	}))
	defer timer.ObserveDuration()

	stream, err := w.Client.Stream(ctx, req, opts...)
	if err != nil {
		status = "failure"
	}
	w.prom.streamCnt.WithLabelValues(w.prom.opts.ServiceID, w.prom.opts.ServiceName, w.prom.opts.ServiceVersion, endpoint, status).Inc()

	return stream, err
}

func (w *wrapper) Publish(ctx context.Context, p client.Message, opts ...client.PublishOption) error {
	status, endpoint := "success", p.Topic()

	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(s float64) {
		w.prom.publishDur.WithLabelValues(w.prom.opts.ServiceID, w.prom.opts.ServiceName, w.prom.opts.ServiceVersion, endpoint).Observe(s)
	}))
	defer timer.ObserveDuration()

	err := w.Client.Publish(ctx, p, opts...)
	if err != nil {
		status = "failure"
	}
	w.prom.publishCnt.WithLabelValues(w.prom.opts.ServiceName, w.prom.opts.ServiceVersion, w.prom.opts.ServiceID, endpoint, status).Inc()

	return err
}

func (w *wrapper) HandlerFunc(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		status, endpoint := "success", req.Endpoint()

		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(s float64) {
			w.prom.handleDur.WithLabelValues(w.prom.opts.ServiceID, w.prom.opts.ServiceName, w.prom.opts.ServiceVersion, endpoint).Observe(s)
		}))
		defer timer.ObserveDuration()

		err := fn(ctx, req, rsp)
		if err != nil {
			status = "failure"
		}
		w.prom.handleCnt.WithLabelValues(w.prom.opts.ServiceID, w.prom.opts.ServiceName, w.prom.opts.ServiceVersion, endpoint, status).Inc()

		return err
	}
}

func (w *wrapper) SubscriberFunc(fn server.SubscriberFunc) server.SubscriberFunc {
	return func(ctx context.Context, msg server.Message) error {
		status, endpoint := "success", msg.Topic()

		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(s float64) {
			w.prom.subscribeDur.WithLabelValues(w.prom.opts.ServiceID, w.prom.opts.ServiceName, w.prom.opts.ServiceVersion, endpoint).Observe(s)
		}))
		defer timer.ObserveDuration()

		err := fn(ctx, msg)
		if err != nil {
			status = "failure"
		}
		w.prom.subscribeCnt.WithLabelValues(w.prom.opts.ServiceID, w.prom.opts.ServiceName, w.prom.opts.ServiceVersion, endpoint, status).Inc()

		return err
	}
}
