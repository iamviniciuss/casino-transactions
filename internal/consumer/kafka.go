// internal/consumer/kafka_consumer.go
package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/iamviniciuss/casino-transactions/internal/use_case"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader   *kafka.Reader
	consumer use_case.TransactionProcessor
}

func NewKafkaConsumer(broker, topic string, uc use_case.TransactionProcessor) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: "transaction-consumer",
	})
	return &KafkaConsumer{reader: reader, consumer: uc}
}

func (kc *KafkaConsumer) Start(ctx context.Context) error {
	for {
		log.Println("Waiting for messages...")
		m, err := kc.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		var tx use_case.ProcessTransactionInput
		if err := json.Unmarshal(m.Value, &tx); err != nil {
			log.Printf("Invalid message: %v - %v", err, string(m.Value))
			continue
		}

		if err := kc.consumer.Process(ctx, tx); err != nil {
			log.Printf("Use case failed: %v", err)
			continue
		}
	}
}
