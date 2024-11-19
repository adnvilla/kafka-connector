package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	kconnector "github.com/adnvilla/kafka-connector"
	"github.com/adnvilla/kafka-connector/base"
)

func main() {
	ctx := context.Background()
	client, err := kconnector.NewClient(base.Config{
		BootstrapServers: []string{"localhost:29092"},
		ClientID:         "example",
		Provider:         base.ZKakfa,
		UseGlobalClient:  true,
	})
	randomNames := []string{"stewy", "lydia", "asif", "mike", "justin"}
	if err != nil {
		log.Panic(err)
	}
	//defer client.Close()
	for {
		event := DummyEvent{
			Name: randomNames[rand.Intn(len(randomNames))],
			Age:  rand.Intn(100),
		}

		err := client.ProduceMessage(ctx, "kafka-connector-example-topic", event)
		if err != nil {
			log.Panic(err)
		}

		time.Sleep(time.Second)
	}
}

// DummyEvent is a deserializable struct for producing/consuming kafka message values.
type DummyEvent struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
