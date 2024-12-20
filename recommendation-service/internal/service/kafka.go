package service

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaProducer(brokers string, topic string) (*KafkaProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers})
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{producer: p, topic: topic}, nil
}

func (kp *KafkaProducer) SendMessage(message interface{}) error {
	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kp.topic, Partition: kafka.PartitionAny},
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
