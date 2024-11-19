package kafka_connector

import (
	"context"
	"github.com/adnvilla/kafka-connector/base"
	"github.com/adnvilla/kafka-connector/zkafka"
)

type KafkaConnector interface {
	ProduceMessage(ctx context.Context, topic string, message interface{}) error
	ConsumeMessages(ctx context.Context, topic, groupId string, handler base.ConsumerHandler) error
	Close() error
}

func NewClient(cfg Config) (KafkaConnector, error) {
	client := &Client{
		cfg: cfg,
	}

	switch cfg.Provider {
	case ZKakfa:
		client.provider = zkafka.GetConnector(base.Config{
			BootstrapServers: cfg.BootstrapServers,
			ClientID:         cfg.ClientID,
			UseGlobalClient:  cfg.UseGlobalClient,
		})
	}

	return client, nil
}

type Client struct {
	cfg      Config
	provider KafkaConnector
}

func (k *Client) ProduceMessage(ctx context.Context, topic string, message interface{}) error {
	return k.provider.ProduceMessage(ctx, topic, message)
}

func (k *Client) ConsumeMessages(ctx context.Context, topic, groupId string, handler base.ConsumerHandler) error {
	return k.provider.ConsumeMessages(ctx, topic, groupId, handler)
}

func (k *Client) Close() error {
	return k.provider.Close()
}
