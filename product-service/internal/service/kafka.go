package service

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type IKafkaProducer interface {
	SendMessage(message interface{}, topic *string) error
	TopicNew() *string
	TopicUpdate() *string
}

type KafkaProducer struct {
	producer    *kafka.Producer
	topicNew    string
	topicUpdate string
}

func (k *KafkaProducer) TopicNew() *string {
	return &k.topicNew
}
func (k *KafkaProducer) TopicUpdate() *string {
	return &k.topicUpdate
}

func NewKafkaProducer(brokers string, topicNew string, topicUpdate string) *KafkaProducer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers})
	if err != nil {
		panic(err)
	}

	return &KafkaProducer{producer: p, topicNew: topicNew, topicUpdate: topicUpdate}
}

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

func (kp *KafkaProducer) Close() error {
	kp.producer.Close()

	return nil
}
