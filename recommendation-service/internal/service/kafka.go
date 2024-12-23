package service

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
	"os"
	"recommendation-service/log"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	logger   log.Factory
	db       IRepo
}

func NewKafkaConsumer(logger log.Factory, repo IRepo) *KafkaConsumer {
	brokers := os.Getenv("KAFKA_BROKER")
	if brokers == "" {
		logger.Bg().Error("KAFKA_BROKER environment variable not set")
		panic("KAFKA_BROKER environment variable not set")
	}

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          "recommendation-service",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		logger.Bg().Error("Failed to create Kafka consumer", zap.Error(err))
		panic(err)
	}

	topics := []string{
		os.Getenv("KAFKA_TOPIC_NEW_USER"),
		os.Getenv("KAFKA_TOPIC_UPDATE_PRODUCT"),
		os.Getenv("KAFKA_TOPIC_NEW_PRODUCT"),
	}

	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		logger.Bg().Error("Failed to subscribe the topics ", zap.Error(err))
		panic(err)
	}
	v, _ := consumer.Subscription()
	for _, t := range v {
		logger.Bg().Info("subscribed topics:", zap.String("topics", t))
	}

	return &KafkaConsumer{
		consumer: consumer,
		logger:   logger,
		db:       repo,
	}
}

func (k *KafkaConsumer) Run(ctx context.Context) error {
	k.logger.Bg().Info("KafkaConsumer is running", zap.String("topic", "product_update"))
	defer k.logger.Bg().Info("KafkaConsumer stopped")

	meta, err := k.consumer.GetMetadata(nil, true, 5000)
	if err != nil {
		k.logger.Bg().Error("Failed to fetch metadata", zap.Error(err))
	} else {
		k.logger.Bg().Info("Consumer metadata", zap.Any("metadata", meta))
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			event := k.consumer.Poll(100)
			if event == nil {
				continue
			}
			switch e := event.(type) {
			case *kafka.Message:
				k.logger.Bg().Info("Received message",
					zap.String("topic", *e.TopicPartition.Topic),
					zap.String("message", string(e.Value)),
					zap.Int32("partition", e.TopicPartition.Partition))
				go func(msg kafka.Message) {
					if err := k.ProcessMessage(msg); err != nil {
						k.logger.Bg().Error("Failed to process message", zap.Error(err))
					}
				}(*e)

			case kafka.Error:
				k.logger.Bg().Error("Kafka error", zap.Error(e))

			default:
				k.logger.Bg().Debug("Ignored event", zap.Any("event", e))
			}
		}
	}
}

func (k *KafkaConsumer) ProcessMessage(msg kafka.Message) error {
	switch *msg.TopicPartition.Topic {
	case "product_update":
		err := k.db.ProductUpdateMsgRepo(msg)
		if err != nil {
			k.logger.Bg().Error("Failed to process message", zap.Error(err))
			return err
		}
	case "product_new":
		err := k.db.ProductNewMsgRepo(msg)
		if err != nil {
			k.logger.Bg().Error("Failed to process message", zap.Error(err))
			return err
		}
	case "user_new":
		err := k.UserNewMsg(msg)
		if err != nil {
			k.logger.Bg().Error("Failed to process message", zap.Error(err))
			return err
		}
	default:
		return fmt.Errorf("unknown topic: %v", *msg.TopicPartition.Topic)
	}
	return nil
}

func (k *KafkaConsumer) ProductUpdateMsgRepo(msg kafka.Message) error {
	return k.db.ProductUpdateMsgRepo(msg)
}

func (k *KafkaConsumer) UserNewMsg(msg kafka.Message) error {
	return k.db.UserNewMsgRepo(msg)
}

func (k *KafkaConsumer) ProductNewMsgRepo(msg kafka.Message) error {
	return k.db.ProductNewMsgRepo(msg)
}

func (k *KafkaConsumer) Stop() error {
	k.logger.Bg().Info("Stopping KafkaConsumer")
	return k.consumer.Close()
}
