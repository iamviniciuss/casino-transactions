package main

import (
	"fmt"

	"github.com/iamviniciuss/casino-transactions/internal/api/http"
	"github.com/iamviniciuss/casino-transactions/internal/api/router"
	"github.com/iamviniciuss/casino-transactions/pkg/config"
)

func main() {
	fmt.Println("Hello, World!")

	http := http.NewFiberHttp()
	router.DataSourceRouter(http)

	configuration := config.NewConfig()

	err := http.ListenAndServe(configuration.Port)
	if err != nil {
		panic(err)
	}
}
