// Description: This is the main entry point for the Casino Transactions Consumer application.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/iamviniciuss/casino-transactions/internal/module/transaction/consumer"
	"github.com/iamviniciuss/casino-transactions/internal/module/transaction/repository"
	"github.com/iamviniciuss/casino-transactions/internal/module/transaction/use_case"
	"github.com/iamviniciuss/casino-transactions/pkg/config"
	"github.com/iamviniciuss/casino-transactions/pkg/shared/message_broker"
)

func main() {
	fmt.Println("Starting Casino Transactions Consumer... 2")

	configuration := config.NewConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConn, err := sql.Open("postgres", configuration.PostgresDSN)
	if err != nil {
		panic(err)
	}

	repo := repository.NewTransactionRepository(dbConn)

	kc := message_broker.NewKafkaConsumer(configuration.KafkaURL, "casino-transactions", "transaction-group")
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
