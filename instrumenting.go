package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	next           KafkaService
}

func (mw instrumentingMiddleware) Consume(ctx context.Context, t string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "consume", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.next.Consume(ctx, t)
	return
}
