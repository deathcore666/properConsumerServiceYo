package main

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   KafkaService
}

func (mw loggingMiddleware) Consume(ctx context.Context, t string) (output string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "consume",
			"input", t,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.next.Consume(ctx, t)
	return
}
