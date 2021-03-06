package kafka

import (
	"fmt"
	"log"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"

	"github.com/Shopify/sarama"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
)

var (
	logger = watermill.NewStdLogger(false, false)
)

func NewSubscriber(consumerGroup string, kafkaHost string, kafkaPort string) *kafka.Subscriber {
	saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
	saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	kakfaBrokerUrl := fmt.Sprintf("%s:%s", kafkaHost, kafkaPort)
	subscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               []string{kakfaBrokerUrl},
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: saramaSubscriberConfig,
			ConsumerGroup:         consumerGroup,
		},
		logger,
	)
	if err != nil {
		log.Panicf("Error creating subscriber Kafka: %+v", err)
	}
	return subscriber
}

func NewPublisher(kafkaHost string, kafkaPort string) *kafka.Publisher {
	kafkaBrokerUrl := fmt.Sprintf("%s:%s", kafkaHost, kafkaPort)
	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   []string{kafkaBrokerUrl},
			Marshaler: kafka.DefaultMarshaler{},
		},
		logger,
	)
	if err != nil {
		log.Panicf("Error creating publisher Kafka: %+v", err)
	}
	return publisher
}

func NewRouter() *message.Router {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	router.AddPlugin(plugin.SignalsHandler)
	router.AddMiddleware(
		middleware.Recoverer,
	)
	return router
}
