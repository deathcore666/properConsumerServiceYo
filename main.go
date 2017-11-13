package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kit/kit/log"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

func main() {
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
		proxy  = flag.String("proxy", "", "Optional comma-separated list of URLs to proxy uppercase requests")
	)
	flag.Parse()

	logger := log.NewLogfmtLogger(os.Stderr)
	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "kafka_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "kafka_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in mciroseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "kafka_service",
		Name:      "count_result",
		Help:      "The result of each method call.",
	}, []string{})

	var svc KafkaService
	svc = kafkaService{}
	svc = proxyingMiddleware(context.Background(), *proxy, logger)(svc)
	svc = loggingMiddleware(logger)(svc)
	svc = instrumentingMiddleware(requestCount, requestLatency, countResult)(svc)

	consumehandler := httptransport.NewServer(
		makeConsumeEndpoint(svc),
		decodeConsumeRequest,
		encodeResponse,
	)

	http.Handle("/consume", consumehandler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", *listen)
	logger.Log("err", http.ListenAndServe(*listen, nil))
}
