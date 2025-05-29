package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/iamviniciuss/casino-transactions/internal/consumer"
	"github.com/iamviniciuss/casino-transactions/internal/repository"
	"github.com/iamviniciuss/casino-transactions/internal/use_case"
	"github.com/iamviniciuss/casino-transactions/pkg/config"
)

func main() {
	configuration := config.NewConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConn, err := sql.Open("postgres", configuration.PostgresDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		panic(err)
	}

	repo := repository.NewTransactionRepository(dbConn)

	kc := consumer.NewKafkaConsumer(configuration.KafkaURL, "casino-transactions", "transaction-group")
	kc.RegisterHandler("transaction", consumer.NewProcessTransactionHandler(use_case.NewProcessTransaction(repo)))

	go func() {
		log.Println("Starting Kafka consumer...")
		if err := kc.Start(ctx); err != nil {
			log.Printf("Kafka consumer error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Received signal %v, terminating gracefully...", sig)

	cancel()

	log.Println("Kafka Consumer terminated successfully")
}
