package config

import (
	"fmt"
	"os"
	"path"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type (
	Config struct {
		Port        string `envconfig:"PORT" default:":9095"`
		KafkaURL    string `envconfig:"KAFKA_URL" default:"kafka:29092"`
		PostgresDSN string `envconfig:"POSTGRES_DSN" default:"host=postgres user=casino password=casino123 dbname=casino_db port=5432 sslmode=disable"`
		// PostgresDSN string `envconfig:"POSTGRES_DSN" default:"host=host.internal.docker user=casino password=casino123 dbname=casino_db port=5432 sslmode=disable"`
	}
)

func NewConfig() Config {
	var c Config
	basePath, _ := os.Getwd()
	godotenv.Load(path.Dir(basePath) + "/.env.local")
	if err := envconfig.Process("", &c); err != nil {
		panic(fmt.Errorf("failed to load config from environment variables: %s", err))
	}
	return c
}
