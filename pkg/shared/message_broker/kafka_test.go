package message_broker

import (
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"

	consumer1 "github.com/iamviniciuss/casino-transactions/internal/module/transaction/consumer"
	"github.com/iamviniciuss/casino-transactions/internal/module/transaction/core"
	"github.com/iamviniciuss/casino-transactions/pkg/test_utils"
	"github.com/segmentio/kafka-go"
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
	ctx := context.Background()

	broker, terminate, err := test_utils.SetupKafkaContainer(t, context.Background())
	require.NoError(t, err)
	defer terminate(t.Context())

	topic := "test-topic"
	groupID := "test-group"
	messageKey := "transaction-key"

	payload := map[string]interface{}{
		"user_id":          "912457bc-7eed-4170-aba5-8a13c35a9d8a",
		"transaction_type": "win",
		"amount":           256.75,
	}
	messageValue, err := json.Marshal(payload)
	require.NoError(t, err)

	require.NoError(t, createKafkaTopic(ctx, broker, topic))

	handler := &mockHandler{}
	consumer := NewKafkaConsumer(broker, topic, groupID)
	consumer.RegisterHandler(messageKey, handler)

	consumerCtx, cancelConsumer := context.WithCancel(context.Background())

	go func() {
		consumer.Start(consumerCtx)
	}()

	err = sendKafkaMessage(broker, topic, messageKey, messageValue)
	require.NoError(t, err)

	waitForConsumerReady()
	cancelConsumer()

	var data consumer1.ProcessTransactionHandlerInput
	err = json.Unmarshal(handler.value, &data)
	require.NoError(t, err)

	assert.True(t, handler.called, "handler should have been called")
	assert.Equal(t, "912457bc-7eed-4170-aba5-8a13c35a9d8a", data.UserID, "user_id should match")
	assert.Equal(t, core.TransactionTypeWin, data.TransactionType, "transaction_type should match")
	assert.Equal(t, 256.75, data.Amount, "amount should match")
}

func createKafkaTopic(ctx context.Context, broker, topic string) error {
	conn, err := kafka.DialContext(ctx, "tcp", broker)
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Println("Creating topic:", topic)
	return conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
}

func sendKafkaMessage(broker, topic, key string, value []byte) error {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	return writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(key),
		Value: value,
	})
}

func waitForConsumerReady() {
	time.Sleep(10 * time.Second)
}
