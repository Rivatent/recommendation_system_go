package app

import (
	"context"
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"os"
	"recommendation-service/internal/closer"
	"recommendation-service/internal/service"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	topics   []string
}

func NewWorker(service *service.Service) (KafkaConsumer, error) {
	topicUser := os.Getenv("KAFKA_TOPIC_USER")
	topicProduct := os.Getenv("KAFKA_TOPIC_PRODUCT")
	if topicUser == "" || topicProduct == "" {
		return KafkaConsumer{}, errors.New("KAFKA_TOPIC_USER or KAFKA_TOPIC_PRODUCT is not set")
	}
	topics := []string{topicUser, topicProduct}
	brokers := os.Getenv("KAFKA_BROKERS")

	kafkaConsumerPtr, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
	})
	if err != nil {
		return KafkaConsumer{}, err
	}
	closer.Add(kafkaConsumerPtr.Close)
	err = kafkaConsumerPtr.SubscribeTopics(topics, nil)
	if err != nil {
		return KafkaConsumer{}, err
	}

	kafkaConsumer := KafkaConsumer{kafkaConsumerPtr, topics}

	return kafkaConsumer, err
}

func (k *KafkaConsumer) Run(ctx context.Context) error {
	for ctx.Err() == nil {
		msg, err := k.consumer.ReadMessage(-1)
		if err == nil {
			log.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
	return ctx.Err()
}
func (k *KafkaConsumer) Stop() error {
	return k.consumer.Close()
}
