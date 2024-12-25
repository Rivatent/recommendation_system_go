package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
	"os"
	"recommendation-service/log"
)

// KafkaConsumer представляет собой консюмер Kafka, который обрабатывает сообщения
// из различных топиков и взаимодействует с репозиторием для обновления данных.
type KafkaConsumer struct {
	consumer *kafka.Consumer
	logger   log.Factory
	db       IRepo
}

// NewKafkaConsumer создает новый экземпляр KafkaConsumer.
// Он считывает настройки, такие как Kafka брокеры и топики, из переменных окружения.
// В случае ошибки при создании консюмера приложение завершает работу.
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

// Run начинает прослушивание Kafka и обработку получаемых сообщений.
// Метод работает в бесконечном цикле, пока контекст не будет отменен.
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

// ProcessMessage обрабатывает входящее сообщение в зависимости от его источника (топика).
// Если сообщение поступило из известного топика, оно передается соответствующему методу обработки.
func (k *KafkaConsumer) ProcessMessage(msg kafka.Message) error {
	topicUpdateProduct := os.Getenv("KAFKA_TOPIC_UPDATE_PRODUCT")
	topicNewProduct := os.Getenv("KAFKA_TOPIC_NEW_PRODUCT")
	topicNewUser := os.Getenv("KAFKA_TOPIC_NEW_USER")

	switch *msg.TopicPartition.Topic {
	case topicUpdateProduct:
		err := k.ProductUpdateMsg(msg)
		if err != nil {
			k.logger.Bg().Error("Failed to process message", zap.Error(err))
			return err
		}
	case topicNewProduct:
		err := k.ProductNewMsg(msg)
		if err != nil {
			k.logger.Bg().Error("Failed to process message", zap.Error(err))
			return err
		}
	case topicNewUser:
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

// ProductUpdateMsg обрабатывает сообщение, связанное с обновлением продукта.
func (k *KafkaConsumer) ProductUpdateMsg(msg kafka.Message) error {
	var updatedProduct map[string]interface{}

	if err := json.Unmarshal(msg.Value, &updatedProduct); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return k.db.ProductUpdateMsgRepo(updatedProduct)
}

// UserNewMsg обрабатывает сообщение о создании нового пользователя.
func (k *KafkaConsumer) UserNewMsg(msg kafka.Message) error {
	var newUser map[string]interface{}

	if err := json.Unmarshal(msg.Value, &newUser); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return k.db.UserNewMsgRepo(newUser)
}

// ProductNewMsg обрабатывает сообщение о создании нового продукта.
func (k *KafkaConsumer) ProductNewMsg(msg kafka.Message) error {
	var newProduct map[string]interface{}

	if err := json.Unmarshal(msg.Value, &newProduct); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return k.db.ProductNewMsgRepo(newProduct)
}

// Stop останавливает консюмер Kafka и закрывает соединение.
func (k *KafkaConsumer) Stop() error {
	k.logger.Bg().Info("Stopping KafkaConsumer")
	return k.consumer.Close()
}
