package main

import (
	"context"
	"fmt"
	"log"

	kconnector "github.com/adnvilla/kafka-connector"
	"github.com/adnvilla/kafka-connector/base"
	"github.com/go-viper/mapstructure/v2"
)

func main() {
	ctx := context.Background()
	client, err := kconnector.NewClient(base.Config{
		BootstrapServers: []string{"localhost:29092"},
		ClientID:         "example",
		Provider:         base.ZKakfa,
		UseGlobalClient:  false,
	})

	if err != nil {
		log.Panic(err)
	}
	// It's important to close the client after consumption to gracefully leave the consumer group
	// (this commits completed work, and informs the broker that this consumer is leaving the group which yields a faster rebalance)
	defer client.Close()

	if err := client.ConsumeMessages(ctx, "kafka-connector-example-topic", "kafka-connector/example/example-consumer",
		func(ctx context.Context, message interface{}) error {

			fmt.Printf("message = %+v\n", message)

			var msg DummyEvent
			mapstructure.Decode(message, &msg)

			fmt.Printf("msg = %+v\n", msg)
			return nil
		},
	); err != nil {
		log.Panic(err)
	}
}

// DummyEvent is a deserializable struct for producing/consuming kafka message values.
type DummyEvent struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
