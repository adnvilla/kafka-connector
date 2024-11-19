package zkafka

import (
	"context"
	"github.com/adnvilla/kafka-connector/base"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/zillow/zfmt"
	"github.com/zillow/zkafka"
)

var (
	conn         Connector
	kafkaOnce    sync.Once
	GetConnector = getConnector
)

type Connector struct {
	cfg      base.Config
	Instance *zkafka.Client
	Error    error
}

func getConnector(cfg base.Config) *Connector {
	if cfg.UseGlobalClient {
		kafkaOnce.Do(func() {
			client := zkafka.NewClient(zkafka.Config{
				BootstrapServers: cfg.BootstrapServers,
			})

			conn.Instance, conn.cfg, conn.Error = client, cfg, nil
		})
	} else {
		client := zkafka.NewClient(zkafka.Config{
			BootstrapServers: cfg.BootstrapServers,
		})

		conn.Instance, conn.cfg, conn.Error = client, cfg, nil
	}

	return &conn
}

func (c *Connector) ProduceMessage(ctx context.Context, topic string, message interface{}) error {
	producer, err := c.Instance.Writer(ctx, zkafka.ProducerTopicConfig{
		ClientID:  c.cfg.ClientID,
		Topic:     topic,
		Formatter: zfmt.JSONFmt,
	})

	if err != nil {
		return err
	}

	_, err = producer.Write(ctx, &message)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connector) ConsumeMessages(ctx context.Context, topic, groupId string, handler base.ConsumerHandler) error {

	topicConfig := zkafka.ConsumerTopicConfig{
		// ClientID is used for caching inside zkafka, and observability within streamz dashboards. But it's not an important
		// part of consumer group semantics. A typical convention is to use the service name executing the kafka worker
		ClientID: c.cfg.ClientID,
		// GroupID is the consumer group. If multiple instances of the same consumer group read messages for the same
		// topic the topic's partitions will be split between the collection. The broker remembers
		// what offset has been committed for a consumer group, and therefore work can be picked up where it was left off
		// across releases
		GroupID: groupId,
		Topic:   topic,
		// The formatter is registered internally to the `zkafka.Message` and used when calling `msg.Decode()`
		// string fmt can be used for both binary and pure strings encoded in the value field of the kafka message. Other options include
		// json, proto, avro, etc.
		Formatter: zfmt.JSONFmt,
		AdditionalProps: map[string]any{
			// only important the first time a consumer group connects. Subsequent connections will start
			// consuming messages
			"auto.offset.reset": "earliest",
		},
	}
	// optionally set up a channel to signal when worker shutdown should occur.
	// A nil channel is also acceptable, but this example demonstrates how to make utility of the signal.
	// The channel should be closed, instead of simply written to, to properly broadcast to the separate worker threads.
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	shutdown := make(chan struct{})

	go func() {
		<-stopCh
		close(shutdown)
	}()

	wf := zkafka.NewWorkFactory(c.Instance)
	// Register a processor which is executed per message.
	// Speedup is used to create multiple processor goroutines. Order is still maintained with this setup by way of `virtual partitions`
	work := wf.Create(topicConfig, &processor{
		handler: handler,
	}, zkafka.Speedup(5))
	if err := work.Run(ctx, shutdown); err != nil {
		return err
	}

	return nil
}

func (c *Connector) Close() error {
	return c.Instance.Close()
}

type processor struct {
	handler base.ConsumerHandler
}

func (p processor) Process(ctx context.Context, msg *zkafka.Message) error {
	var event interface{}
	err := msg.Decode(&event)
	if err != nil {
		return err
	}

	// optionally, if you don't want to use the configured formatter at all, access the kafka message payload bytes directly.
	// The commented out block shows accessing the byte array. In this case we're stringifying the bytes, but this could be json unmarshalled,
	// proto unmarshalled etc., depending on the expected payload
	//data := msg.Value()
	//str := string(data)

	err = p.handler(ctx, event)
	if err != nil {
		return err
	}

	//log.Printf(" offset: %d, partition: %d. event.Name: %s, event.Age %d\n", msg.Offset, msg.Partition, event.Name, event.Age)
	return nil
}

func ResetConnection() {
	kafkaOnce = sync.Once{}
}
