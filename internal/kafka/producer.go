package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/viper"
	"github.com/aslbnv/messagio/internal/types"
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
	msgStr := fmt.Sprintf("[id]: %s [text]: %s [created_at]: %s\n", msg.ID, msg.Text, msg.CreatedAt)
	topic := viper.GetString("kafka.topic")

	return k.producer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(msgStr),
		}, nil,
	)
}
