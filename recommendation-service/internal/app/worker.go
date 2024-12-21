package app

import (
	"context"
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
	"os"
	"recommendation-service/internal/service"
	"recommendation-service/log"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	logger   log.Factory
	service  *service.Service
}

func NewKafkaConsumer(logger log.Factory, svc *service.Service) (*KafkaConsumer, error) {
	brokers := os.Getenv("KAFKA_BROKER")
	if brokers == "" {
		logger.Bg().Error("KAFKA_BROKER environment variable not set")
		return nil, errors.New("KAFKA_BROKER environment variable not set")
	}

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          "recommendation-service",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		logger.Bg().Error("Failed to create Kafka consumer", zap.Error(err))
		return nil, err
	}

	err = consumer.Subscribe("product_updates", nil)
	if err != nil {
		logger.Bg().Error("Failed to subscribe to topic", zap.String("topic", "product_updates"), zap.Error(err))
		return nil, err
	}

	return &KafkaConsumer{
		consumer: consumer,
		logger:   logger,
		service:  svc,
	}, nil
}

func (k *KafkaConsumer) Run(ctx context.Context) error {
	k.logger.Bg().Info("KafkaConsumer is running", zap.String("topic", "product_updates"))
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
			event := k.consumer.Poll(100) // Ожидание событий с timeout 100 мс
			if event == nil {
				continue
			}

			switch e := event.(type) {
			case *kafka.Message:
				k.logger.Bg().Info("Received message",
					zap.String("topic", *e.TopicPartition.Topic),
					zap.String("message", string(e.Value)),
					zap.Int32("partition", e.TopicPartition.Partition))
				//zap.Int64("offset", e.TopicPartition.Offset))

				// Обработка сообщения
				// err := k.service.ProcessMessage(e.Value)
				// if err != nil {
				//     k.logger.Bg().Error("Failed to process message", zap.Error(err))
				// }

			case kafka.Error:
				k.logger.Bg().Error("Kafka error", zap.Error(e))

			default:
				k.logger.Bg().Debug("Ignored event", zap.Any("event", e))
			}
		}
	}
}

func (k *KafkaConsumer) Stop() error {
	k.logger.Bg().Info("Stopping KafkaConsumer")
	return k.consumer.Close()
}
