package main

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

type logmw struct {
	logger log.Logger
	next   KafkaService
}

func loggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next KafkaService) KafkaService {
		return logmw{logger, next}
	}
}

func (mw logmw) Consume(ctx context.Context, t string) (output string, err error) {
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
