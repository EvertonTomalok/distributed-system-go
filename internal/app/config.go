package app

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/evertontomalok/distributed-system-go/internal/app/shared"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/dto"
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
	Mongodb struct {
		Host string
	}
	UserApi struct {
		BaseUrl string
	}
}

func Configure() Config {
	const LocalHost = "0.0.0.0"

	viper.SetDefault("Host", LocalHost)
	viper.SetDefault("Port", "5000")
	viper.SetDefault("Postgres.Host", "postgres://postgres:secret@db:5432/distributed-system?sslmode=disable")
	viper.SetDefault("Kafka.Host", "kafka")
	viper.SetDefault("Kafka.Port", "29092")
	viper.SetDefault("Mongodb.Host", "mongodb://root:secret@mongodb:27017/?maxPoolSize=20&w=majority")
	viper.SetDefault("UserApi.BaseUrl", "http://0.0.0.0:8000/api/")

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

	if err := postgres.Check(database); err != nil {
		log.Fatal("Database not initilialized.")
	}

	adapter := postgres.New(database)
	// Dependency Injections
	methods.MethodsDBAdapter = adapter
	orders.OrdersDBAdapter = adapter

	log.Infof("Database connection is ready at [%s***:%s]", cfg.Postgres.Host[0:2], cfg.Port)
}

func InitKafka(ctx context.Context, cfg Config) {
	kafka.Publisher = kafka.NewPublisher(cfg.Kafka.Host, cfg.Kafka.Port)
	log.Infof("Kafka is ready at [%s***:%s]", cfg.Kafka.Host[0:2], cfg.Kafka.Port)
}

func ConfigureFlags() {
	// Todo implement some logic simulating getting values from a cache service with real time watch changes,
	//	like redis os etcd
	shared.Flags = map[string]bool{
		dto.PostOrderFlag:         true,
		dto.GetOrderByIdFlag:      true,
		dto.GetOrdersFromUserFlag: true,
	}

	log.Info("Feature Flags configured.")
}
