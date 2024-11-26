package kafka

import (
	"context"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

// Kafka consumer struct
type KafkaConsumer struct {
	brokers []string
	topic   string
	group   string
	logger  *slog.Logger
}

// Kafka consumer constructor
func NewKafkaConsumer(brokers []string, topic, group string, logger *slog.Logger) *KafkaConsumer {
	return &KafkaConsumer{
		brokers: brokers,
		topic:   topic,
		group:   group,
		logger:  logger,
	}
}

// Create kafka consumer
func (kc *KafkaConsumer) Consume(handler func([]byte) error) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  kc.brokers,
		GroupID:  kc.group,
		Topic:    kc.topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	kc.logger.Info("kafka consumer is ready to consume", "topic", kc.topic)

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			kc.logger.Error("error while reading message", "error", err)
			return err
		}

		kc.logger.Info("consumed message", "message", string(m.Value))

		err = handler(m.Value)
		if err != nil {
			kc.logger.Error("error while processing message", "error", err)
		} else {
			kc.logger.Info("successfully processed message")
		}
	}
}
