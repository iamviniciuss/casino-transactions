package message_broker

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaHandler interface {
	Handle(ctx context.Context, msg []byte) error
}

type KafkaConsumer struct {
	reader   *kafka.Reader
	handlers map[string]KafkaHandler
}

func NewKafkaConsumer(broker, topic, groupID string) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:       []string{broker},
		Topic:         topic,
		GroupID:       groupID,
		QueueCapacity: 100,
		MaxWait:       10 * 1000000,
	})

	return &KafkaConsumer{
		reader:   reader,
		handlers: make(map[string]KafkaHandler),
	}
}

func (kc *KafkaConsumer) RegisterHandler(key string, handler KafkaHandler) {
	kc.handlers[key] = handler
}

func (kc *KafkaConsumer) Start(ctx context.Context) error {
	for {
		m, err := kc.reader.ReadMessage(ctx)
		log.Println("Received message:", string(m.Value), "with key:", string(m.Key))
		if err != nil {
			return err
		}

		handlerKey := m.Key
		handler, ok := kc.handlers[string(handlerKey)]
		if !ok {
			log.Printf("No handler for key: %s", string(handlerKey))
			continue
		}

		if err := handler.Handle(ctx, m.Value); err != nil {
			log.Printf("Handler failed: %v", err)
			continue
		}
	}
}
