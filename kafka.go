package main

import (
	"fmt"

	"github.com/Shopify/sarama"
)

var (
	brokers = []string{"localhost:9092"}
)

func kafkaRoutine(inChan chan string, topic string) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(err)
	}

	topics, _ := consumer.Topics()
	if !(containsTopic(topics, topic)) {
		inChan <- "There is no such a topic"
		fmt.Println("kafkaroutine exited")
		return
	}

	partitionList, err := consumer.Partitions(topic)
	for _, partition := range partitionList {
		pc, _ := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
		go func(pc sarama.PartitionConsumer) {
			for {
				select {
				case msg := <-pc.Messages():
					inChan <- string(msg.Value)
				}
			}
		}(pc)
	}
	fmt.Println("kafkaRoutine exited")
}

func containsTopic(topics []string, topic string) bool {
	for _, v := range topics {
		if topic == v {
			return true
		}
	}
	return false
}
