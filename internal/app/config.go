package app

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/evertontomalok/distributed-system-go/internal/domain/methods"
	"github.com/evertontomalok/distributed-system-go/internal/domain/orders"
	"github.com/evertontomalok/distributed-system-go/internal/infra/kafka"
	"github.com/evertontomalok/distributed-system-go/internal/infra/postgres"
	"github.com/spf13/viper"
)

type Config struct {
	Port     string
	Host     string
	Postgres struct {
		Host string
	}
	Kafka struct {
		Host string
		Port string
	}
}

func Configure() Config {
	const LocalHost = "0.0.0.0"

	viper.SetDefault("Host", LocalHost)
	viper.SetDefault("Port", "5000")
	viper.SetDefault("Postgres.Host", "postgres://postgres:secret@127.0.0.1:5432/distributed-system?sslmode=disable")
	viper.SetDefault("Kafka.Host", LocalHost)
	viper.SetDefault("Kafka.Port", "9092")

	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Errorf("It was impossible configure Server. %+v", err)
	}

	log.Info("Configuration is ready!")

	return cfg
}

func InitDB(ctx context.Context, cfg Config) {
	database := postgres.Init(ctx, cfg.Postgres.Host)
	adapter := postgres.New(database)

	// Dependency Injections
	methods.MethodsDBAdapter = adapter
	orders.OrdersDBAdapter = adapter
	log.Info("Database connection is ready!")
}

func InitKafka(ctx context.Context, cfg Config) {
	kafka.Publisher = kafka.NewPublisher(cfg.Kafka.Host, cfg.Kafka.Port)
	log.Info("Kafka is ready!")
}
