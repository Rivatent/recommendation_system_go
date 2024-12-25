package service

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// IKafkaProducer - интерфейс для Kafka продюсера.
// Определяет методы, которые должен реализовать любой Kafka продюсер.
type IKafkaProducer interface {
	SendMessage(message interface{}, topic *string) error
	TopicNew() *string
	TopicUpdate() *string
}

// KafkaProducer - структура Kafka продюсера.
// Содержит продюсер Kafka и названия топиков.
type KafkaProducer struct {
	producer    *kafka.Producer
	topicNew    string
	topicUpdate string
}

// TopicNew возвращает указатель на название темы для новых сообщений.
func (k *KafkaProducer) TopicNew() *string {
	return &k.topicNew
}

// TopicUpdate возвращает указатель на название темы для обновленных сообщений.
func (k *KafkaProducer) TopicUpdate() *string {
	return &k.topicUpdate
}

// NewKafkaProducer создает новый экземпляр KafkaProducer.
// Принимает строку brokers для подключения к серверам Kafka,
// названия тем для новых и обновленных сообщений.
// Возвращает указатель на созданный экземпляр KafkaProducer.
// В случае ошибки при создании продюсера вызывает panic.
func NewKafkaProducer(brokers string, topicNew string, topicUpdate string) *KafkaProducer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers})
	if err != nil {
		panic(err)
	}

	return &KafkaProducer{producer: p, topicNew: topicNew, topicUpdate: topicUpdate}
}

// SendMessage отправляет сообщение в заданную тему Kafka.
// Принимает сообщение в виде интерфейса и указатель на тему,
// в которую нужно отправить сообщение.
// Сообщение сериализуется в JSON формат.
// В случае ошибки возвращает её.
func (kp *KafkaProducer) SendMessage(message interface{}, topic *string) error {
	msg, err := json.Marshal(message)

	if err != nil {
		return err
	}
	err = kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: topic, Partition: kafka.PartitionAny},
		Value:          msg,
	}, nil)
	if err != nil {
		return err
	}
	kp.producer.Flush(1000)

	return nil
}

// Close закрывает продюсер и освобождает ресурсы.
func (kp *KafkaProducer) Close() error {
	kp.producer.Close()

	return nil
}
