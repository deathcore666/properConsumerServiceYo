package main

import (
	"context"
	"errors"
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
	/* TODO: implement kafka logic here*/
	return "Messages...", nil
}

//ServiceMiddleware is a chainable thing for the service
type ServiceMiddleware func(KafkaService) KafkaService
