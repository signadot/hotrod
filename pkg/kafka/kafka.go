package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"github.com/signadot/hotrod/pkg/config"
	"go.opentelemetry.io/otel/trace"
)

const (
	DispatchDriverTopic = "DispatchDriver"
)

func GetSyncProducer(clientID string, provider trace.TracerProvider) (sarama.SyncProducer, error) {
	conf := getConfig(clientID)
	conf.Producer.Partitioner = sarama.NewRandomPartitioner
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Retry.Max = 10 // Retry up to 10 times to produce the message
	conf.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(config.GetKafkaBrokers(), conf)
	if err != nil {
		return nil, err
	}

	producer = otelsarama.WrapSyncProducer(conf,
		producer,
		otelsarama.WithTracerProvider(provider),
	)
	return producer, nil
}

func GetConsumerGroup(clientID, groupID string, provider trace.TracerProvider,
	handler sarama.ConsumerGroupHandler) (sarama.ConsumerGroup, sarama.ConsumerGroupHandler, error) {
	conf := getConfig(clientID)

	consumerGroup, err := sarama.NewConsumerGroup(
		config.GetKafkaBrokers(),
		config.SignadotConsumerGroup(groupID),
		conf)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating consumer group client: %v", err)
	}

	instrumentedHandler := otelsarama.WrapConsumerGroupHandler(
		handler,
		otelsarama.WithTracerProvider(provider),
	)
	return consumerGroup, instrumentedHandler, nil
}

func getConfig(clientID string) *sarama.Config {
	conf := sarama.NewConfig()
	conf.ClientID = clientID
	conf.Version = sarama.V1_1_0_0
	conf.Net.TLS.Enable = false
	conf.Net.SASL.Enable = false
	return conf
}
