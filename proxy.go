package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/go-kit/kit/sd/lb"

	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"

	httptransport "github.com/go-kit/kit/transport/http"
)

func proxyingMiddleware(ctx context.Context, instances string, logger log.Logger) ServiceMiddleware {
	if instances == "" {
		logger.Log("proxy_to", "none")
		return func(next KafkaService) KafkaService { return next }
	}

	var (
		qps         = 100
		maxAttempts = 3
		maxTime     = 250 * time.Millisecond
	)

	var (
		instanceList = split(instances)
		endpointer   sd.FixedEndpointer
	)
	logger.Log("proxy_to", fmt.Sprint(instanceList))
	for _, instance := range instanceList {
		var e endpoint.Endpoint
		e = makeConsumeProxy(ctx, instance)
		e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
		e = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), qps))(e)
		endpointer = append(endpointer, e)
	}

	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(maxAttempts, maxTime, balancer)

	return func(next KafkaService) KafkaService {
		return proxymw{ctx, retry}
	}
}

func (mw proxymw) Consume(_ context.Context, t string) (string, error) {
	response, err := mw.consume(mw.ctx, consumeRequest{T: t})
	if err != nil {
		return "", err
	}

	resp := response.(consumeResponse)
	if resp.Err != "" {
		return resp.M, errors.New(resp.Err)
	}
	return resp.M, nil

}

func makeConsumeProxy(ctx context.Context, instance string) endpoint.Endpoint {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		panic(err)
	}
	if u.Path == "" {
		u.Path = "/consume"
	}
	return httptransport.NewClient(
		"GET",
		u,
		encodeRequest,
		decodeConsumeResponse,
	).Endpoint()
}

func split(s string) []string {
	a := strings.Split(s, ",")
	for i := range a {
		a[i] = strings.TrimSpace(a[i])
	}
	return a
}

type proxymw struct {
	ctx     context.Context
	consume endpoint.Endpoint
}
