package service

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// IKafkaProducer - интерфейс для Kafka продюсера.
// Определяет методы для отправки сообщений в Kafka и получения имен топиков.
type IKafkaProducer interface {
	SendMessage(message interface{}, topic *string) error
	TopicNew() *string
	TopicUpdate() *string
}

// KafkaProducer - структура, реализующая интерфейс IKafkaProducer.
// Содержит Kafka продюсер и имена топиков для новых и обновлений сообщений.
type KafkaProducer struct {
	producer    *kafka.Producer
	topicNew    string
	topicUpdate string
}

// TopicNew возвращает указатель на имя топика для новых пользователей.
func (k *KafkaProducer) TopicNew() *string {
	return &k.topicNew
}

// TopicUpdate возвращает указатель на имя топика для обновлений пользователей.
func (k *KafkaProducer) TopicUpdate() *string {
	return &k.topicUpdate
}

// NewKafkaProducer создает новый экземпляр KafkaProducer.
// brokers - список адресов Kafka брокеров, используемых для подключения.
// topicNew - имя топика для новых сообщений.
// topicUpdate - имя топика для обновлений сообщений.
// Возвращает указатель на новый экземпляр KafkaProducer.
func NewKafkaProducer(brokers string, topicNew string, topicUpdate string) *KafkaProducer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers})
	if err != nil {
		panic(err)
	}

	return &KafkaProducer{producer: p, topicNew: topicNew, topicUpdate: topicUpdate}
}

// SendMessage отправляет сообщение в указанный топик.
// message - объект, который будет сериализован в JSON.
// topic - указатель на строку с именем топика, в который должно быть отправлено сообщение.
// Возвращает ошибку, если во время отправки сообщения произошла ошибка.
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

// Close закрывает продюсер и освобождает все связанные ресурсы.
// Возвращает ошибку, если закрытие прошло неудачно.
func (kp *KafkaProducer) Close() error {
	kp.producer.Close()
	return nil
}
