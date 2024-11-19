package kafka_connector

import (
	"context"
	"errors"
	"fmt"
	"github.com/adnvilla/kafka-connector/base"
	"github.com/adnvilla/kafka-connector/zkafka"
)

type KafkaConnector interface {
	ProduceMessage(ctx context.Context, topic string, message interface{}) error
	ConsumeMessages(ctx context.Context, topic, groupId string, handler base.ConsumerHandler) error
	Close() error
}

func NewClient(cfg base.Config) (KafkaConnector, error) {
	client := &Client{
		cfg: cfg,
	}

	switch cfg.Provider {
	case base.ZKakfa:
		client.provider = zkafka.GetConnector(base.Config{
			BootstrapServers: cfg.BootstrapServers,
			ClientID:         cfg.ClientID,
			UseGlobalClient:  cfg.UseGlobalClient,
		})
	default:
		return nil, errors.New(fmt.Sprintf("Not supported provider: %v", cfg.Provider))
	}

	return client, nil
}

type Client struct {
	cfg      base.Config
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
