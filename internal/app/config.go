package app

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/evertontomalok/distributed-system-go/internal/domain/methods"
	"github.com/evertontomalok/distributed-system-go/internal/domain/orders"
	"github.com/evertontomalok/distributed-system-go/internal/infra/postgres"
	"github.com/spf13/viper"
)

type Config struct {
	Port     string
	Host     string
	Postgres struct {
		Host string
	}
}

func Configure() Config {
	const LocalHost = "0.0.0.0"

	viper.SetDefault("Host", LocalHost)
	viper.SetDefault("Port", "5000")
	viper.SetDefault("Postgres.Host", "postgres://postgres:secret@127.0.0.1:5432/distributed-system?sslmode=disable")

	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Errorf("It was impossible configure Server. %+v", err)
	}

	return cfg
}

func InitDB(ctx context.Context, cfg Config) {
	database := postgres.Init(ctx, cfg.Postgres.Host)
	adapter := postgres.New(database)
	methods.MethodsDBAdapter = adapter
	orders.OrdersDBAdapter = adapter
}
