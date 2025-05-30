package consumer

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"

	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockHandler struct {
	called bool
	value  []byte
}

func (m *mockHandler) Handle(ctx context.Context, msg []byte) error {
	m.called = true
	m.value = msg
	return nil
}

func TestKafkaConsumer_Start(t *testing.T) {
	// Setup Kafka container
	// broker, terminate := test_utils.SetupKafka(t)
	// defer terminate()

	topic := "test-topic"
	groupID := "test-group"
	messageKey := "my-key"
	messageValue := []byte("hello world")

	broker := "localhost:9092"
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	conn, err := kafka.DialContext(ctx, "tcp", broker)
	require.NoError(t, err)
	defer conn.Close()

	for i := 0; i < 10; i++ {
		log.Println("Creating topic:", topic)
		err = conn.CreateTopics(kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	require.NoError(t, err)

	handler := &mockHandler{}

	kc := NewKafkaConsumer(broker, topic, groupID)
	kc.RegisterHandler(messageKey, handler)

	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()

	go func() {
		err := kc.Start(ctx2)
		if err != nil {
			t.Logf("consumer stopped: %v", err)
		}
	}()

	time.Sleep(4 * time.Second)

	writer := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(messageKey),
		Value: messageValue,
	})
	require.NoError(t, err)

	time.Sleep(10 * time.Second)

	assert.True(t, handler.called, "handler should have been called")
	assert.Equal(t, messageValue, handler.value, "message value should match")
}
