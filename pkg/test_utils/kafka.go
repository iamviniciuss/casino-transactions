package test_utils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupKafkaContainer(t *testing.T, ctx context.Context) (string, func(context.Context) error, error) {
	kafkaContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "bitnami/kafka:3.6",
			ExposedPorts: []string{"9094:9094/tcp"},
			Env: map[string]string{
				"KAFKA_CFG_NODE_ID":                        "0",
				"KAFKA_CFG_PROCESS_ROLES":                  "controller,broker",
				"KAFKA_CFG_CONTROLLER_QUORUM_VOTERS":       "0@localhost:9093",
				"KAFKA_CFG_LISTENERS":                      "PLAINTEXT://:9092,CONTROLLER://:9093,EXTERNAL://:9094",
				"KAFKA_CFG_ADVERTISED_LISTENERS":           "PLAINTEXT://localhost:9092,EXTERNAL://localhost:9094",
				"KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP": "CONTROLLER:PLAINTEXT,EXTERNAL:PLAINTEXT,PLAINTEXT:PLAINTEXT",
				"KAFKA_CFG_CONTROLLER_LISTENER_NAMES":      "CONTROLLER",
				"KAFKA_CFG_INTER_BROKER_LISTENER_NAME":     "PLAINTEXT",
				"ALLOW_PLAINTEXT_LISTENER":                 "yes",
			},
			WaitingFor: wait.ForLog("Kafka Server started").WithStartupTimeout(90 * time.Second),
		},
		Started: true,
	})
	require.NoError(t, err)

	host, err := kafkaContainer.Host(ctx)
	require.NoError(t, err)

	port, err := kafkaContainer.MappedPort(ctx, "9094")
	require.NoError(t, err)

	broker := fmt.Sprintf("%s:%s", host, port.Port())

	return broker, func(ctx context.Context) error {
		return kafkaContainer.Terminate(ctx)
	}, nil
}
