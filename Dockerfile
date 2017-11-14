FROM golang

ADD . /go/src/github.com/deathcore666/ProperConsumerServiceYo

RUN go get github.com/go-kit/kit/metrics/prometheus && go get github.com/sony/gobreaker && go get golang.org/x/time/rate
RUN go get github.com/Shopify/sarama && go get github.com/prometheus/client_golang/prometheus
RUN go get github.com/go-kit/kit/transport/http && go get github.com/go-kit/kit/log && go get github.com/prometheus/client_golang/prometheus/promhttp
RUN go get github.com/go-kit/kit/endpoint && go get github.com/afex/hystrix-go/hystrix && go get github.com/streadway/handy/breaker

RUN go install github.com/deathcore666/ProperConsumerServiceYo

ENTRYPOINT /go/bin/ProperConsumerServiceYo

EXPOSE 8080