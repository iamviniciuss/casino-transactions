// Description: Main entry point for the Casino Transactions API service.
package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/iamviniciuss/casino-transactions/internal/api/http"
	"github.com/iamviniciuss/casino-transactions/internal/api/router"
	"github.com/iamviniciuss/casino-transactions/internal/repository"
	"github.com/iamviniciuss/casino-transactions/pkg/config"
)

func main() {
	fmt.Println("Starting Casino Transactions API...")

	configuration := config.NewConfig()

	dbConn, err := sql.Open("postgres", configuration.PostgresDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	repo := repository.NewTransactionRepository(dbConn)

	http := http.NewFiberHttp()
	router.DataSourceRouter(http, repo)

	err = http.ListenAndServe(configuration.Port)
	if err != nil {
		panic(err)
	}
}
