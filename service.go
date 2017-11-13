package main

import (
	"context"
	"errors"
	"time"
)

//KafkaService yolo
type KafkaService interface {
	Consume(context.Context, string) (string, error)
}

//ErrEmpty yolo
var ErrEmpty = errors.New("No topic provided")

type kafkaService struct{}

//Consumer logic implemented here
func (kafkaService) Consume(_ context.Context, topic string) (string, error) {
	if topic == "" {
		return "", ErrEmpty
	}

	var inChan = make(chan string)
	var readyChan = make(chan struct{})
	var result string
	go kafkaRoutine(inChan, topic)
	go func() {
		for {
			select {
			case msg := <-inChan:
				result = result + msg + "\n"
			case <-time.After(time.Second * 1):
				readyChan <- struct{}{}
			}

		}
	}()

	<-readyChan

	return result, nil

}

//ServiceMiddleware is a chainable thing for the service
type ServiceMiddleware func(KafkaService) KafkaService
