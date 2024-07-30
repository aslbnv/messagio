package kafka

import (
	"fmt"

	"github.com/aslbnv/messagio/internal/types"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/viper"
)

type KafkaProducer struct {
	producer *kafka.Producer
}

func NewKafkaProducer() (*KafkaProducer, error) {
	var (
		host = viper.GetString("kafka.host")
		port = viper.GetString("kafka.port")
	)

	producerConfig := &kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", host, port),
	}

	producer, err := kafka.NewProducer(producerConfig)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		producer: producer,
	}, nil
}

func (k *KafkaProducer) ProduceMessage(msg *types.Message) error {
	msgStr := fmt.Sprintf("[id]: %s [uuid]: %s [text]: %s [created_at]: %s\n", string(msg.ID), msg.UUID, msg.Text, msg.CreatedAt)
	topic := viper.GetString("kafka.topic")

	return k.producer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(msgStr),
		}, nil,
	)
}
