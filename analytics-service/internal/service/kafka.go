package service

import (
	"analytics-service/pkg/log"
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
	"os"
)

// KafkaConsumer - структура, представляющая потребителя Kafka.
// Она включает ссылки на Kafka consumer, логгер и интерфейс репозитория для работы с данными.
type KafkaConsumer struct {
	consumer *kafka.Consumer
	logger   log.Factory
	db       IRepo
}

// NewKafkaConsumer - функция для создания нового экземпляра KafkaConsumer.
// Она настраивает Kafka потребителя, подписывается на заданные темы и возвращает указатель на новый KafkaConsumer.
//
// Параметры:
// - logger log.Factory: логгер, который будет использован для записи событий и ошибок.
// - repo IRepo: интерфейс репозитория для взаимодействия с данными.
//
// Возвращаемое значение:
// - *KafkaConsumer: указатель на структуру KafkaConsumer, готовую к работе.
//
// Примечание: В случае отсутствия переменной окружения "KAFKA_BROKER" или ошибок на этапе
// создания потребителя функция вызывает панику.
func NewKafkaConsumer(logger log.Factory, repo IRepo) *KafkaConsumer {
	brokers := os.Getenv("KAFKA_BROKER")
	if brokers == "" {
		logger.Bg().Error("KAFKA_BROKER environment variable not set")
		panic("KAFKA_BROKER environment variable not set")
	}

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          "analytics-service",
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

// Run - метод, запускающий циклическую обработку сообщений из Kafka.
// Он слушает события и обрабатывает входящие сообщения в зависимости от их типа.
// Метод также учитывает контекст, чтобы корректно завершить работу по его сигналу.
//
// Параметры:
// - ctx context.Context: контекст для контроля завершения работы.
//
// Возвращаемое значение:
// - error: nil, если все прошло успешно; иначе ошибка, произошедшая при обработке.
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

// ProcessMessage - метод для обработки входящих сообщений из Kafka.
// Он определяет, к какой теме относится сообщение, и вызывает соответствующий метод
// для обновления аналитики, если оно принадлежит известной теме.
//
// Параметры:
// - msg kafka.Message: входящее сообщение из Kafka, которое требуется обработать.
//
// Возвращаемое значение:
// - error: nil, если сообщение было успешно обработано; в противном случае возвращается ошибка.
func (k *KafkaConsumer) ProcessMessage(msg kafka.Message) error {
	topicUpdateProduct := os.Getenv("KAFKA_TOPIC_UPDATE_PRODUCT")
	topicNewProduct := os.Getenv("KAFKA_TOPIC_NEW_PRODUCT")
	topicNewUser := os.Getenv("KAFKA_TOPIC_NEW_USER")

	switch *msg.TopicPartition.Topic {
	case topicUpdateProduct:
		fallthrough
	case topicNewProduct:
		fallthrough
	case topicNewUser:
		err := k.UpdateAnalyticsMsg()
		if err != nil {
			k.logger.Bg().Error("Failed to process message", zap.Error(err))
			return err
		}
	default:
		return fmt.Errorf("unknown topic: %v", *msg.TopicPartition.Topic)
	}
	return nil
}

// UpdateAnalyticsMsg - метод для обновления аналитических данных в репозитории.
// Этот метод вызывает функцию репозитория, которая отвечает за изменение данных.
//
// Возвращаемое значение:
// - error: nil, если обновление прошло успешно; иначе ошибка, произошедшая при обращении в репозиторий.
func (k *KafkaConsumer) UpdateAnalyticsMsg() error {

	return k.db.UpdateAnalyticsMsgRepo()
}

// Stop - метод для остановки KafkaConsumer и закрытия соединения с Kafka.
// Он записывает информацию о завершении работы потребителя и завершает его.
//
// Возвращаемое значение:
// - error: nil, если остановка прошла успешно; иначе ошибка, возникшая при закрытии соединения.
func (k *KafkaConsumer) Stop() error {
	k.logger.Bg().Info("Stopping KafkaConsumer")
	return k.consumer.Close()
}
