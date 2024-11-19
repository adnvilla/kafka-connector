# kafka-connector

kafka-conector is a flexible Go library designed to facilitate seamless integration with different Kafka client implementations

[![License](https://img.shields.io/github/license/adnvilla/kafka-connector)](https://github.com/adnvilla/kafka-connector/blob/main/LICENSE)
[![GitHub Actions](https://github.com/adnvilla/kafka-connector/actions/workflows/go.yml/badge.svg)](https://github.com/adnvilla/kafka-connector/actions/workflows/go.yml)
[![Codecov](https://codecov.io/gh/adnvilla/kafka-connector/branch/main/graph/badge.svg?token=STRT8T67YP)](https://codecov.io/gh/adnvilla/kafka-connector)
[![Go Report Card](https://goreportcard.com/badge/github.com/adnvilla/kafka-connector)](https://goreportcard.com/report/github.com/adnvilla/kafka-connector)

## Overview

`Kafka-conector` is a flexible Go library designed to facilitate seamless integration with different Kafka client implementations. By abstracting the underlying Kafka libraries, `kafka-conector` allows developers to easily switch between popular Kafka client libraries without changing application logic. This enables a more adaptable and maintainable system when working with Kafka, offering the flexibility to choose the best-suited library for your use case or scale. Whether you're using sarama, confluent-kafka-go, or another client, kafka-conector ensures that your code remains decoupled and efficient.

---

## Features

- **KafkaConnector Interface**:
  - `ProduceMessage`: Produces messages to a specified topic.
  - `ConsumeMessages`: Consumes messages from a topic in a consumer group and executes a defined handler.
  - `Close`: Closes the Kafka connection.

- **Support for Multiple Providers**:
  - Initial implementation based on **ZKafka** (Zillow Kafka Client).

- **Kafka Client Management**:
  - Use of a global client or creation of independent clients based on configuration.

- **Optimized Message Consumption**:
  - `ConsumerTopicConfig` setup with options like `auto.offset.reset` and safe goroutine handling.
  - Parallelism options with `zkafka.Speedup`.

- **Message Production**:
  - JSON message formatting using `zfmt.JSONFmt`.

- **Flexible Configuration Interface**:
  - Customization through the `base.Config` struct.

---

## Highlights

- **Parallelism**: Support for concurrent message processing via speedup configurations (`zkafka.Speedup`).
- **Extensibility**: Architecture designed to support multiple providers in the future.
- **Simplicity**: Clear and straightforward API for producing and consuming Kafka messages.
- **ZKafka Integration**: Leverages `zkafka` client for robustness and flexibility.

---

## Usage Example

```go
package main

import (
	"context"
	"fmt"
	"github.com/adnvilla/kafka-connector/base"
	"github.com/adnvilla/kafka-connector/kafka_connector"
)

func main() {
	// Configuration
	cfg := base.Config{
		BootstrapServers: []string{"localhost:9092"},
		ClientID:         "example-client",
		Provider:         base.ZKakfa,
		UseGlobalClient:  true,
	}

	client, err := kafka_connector.NewClient(cfg)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}
	defer client.Close()

	// Producing messages
	err = client.ProduceMessage(context.Background(), "example-topic", map[string]string{"key": "value"})
	if err != nil {
		fmt.Printf("Error producing message: %v\n", err)
		return
	}

	// Consuming messages
	err = client.ConsumeMessages(context.Background(), "example-topic", "example-group", func(ctx context.Context, message interface{}) error {
		fmt.Printf("Message received: %v\n", message)
		return nil
	})
	if err != nil {
		fmt.Printf("Error consuming messages: %v\n", err)
	}
}
```

---

## Known Issues

1. **Exclusive Use of ZKafka**:
   - Currently, only the `zkafka` provider is implemented.
   - Other Kafka providers are not yet supported.

2. **Global Client Handling**:
   - If a global client (`UseGlobalClient: true`) is used, changing the configuration requires restarting the application.
   - The `ResetConnection` function allows resetting the client, but it may not be suitable for production use.

3. **Error Handling**:
   - Exceptions in consumer handlers are not directly logged.

4. **Parallelism**:
   - While `zkafka.Speedup` enables parallelism, fine-tuning may be required to avoid overloading the system.

---

## Proposals for Improvement

1. **Support for More Providers**:
   - Add support for other Kafka clients like `confluent-kafka-go` or the official Kafka client for Go.

2. **Metrics and Logging**:
   - Incorporate monitoring tools (e.g., Prometheus, OpenTelemetry) for performance visibility.

3. **Improved Error Handling**:
   - Implement detailed logging for errors in message production and consumption.

4. **Documentation**:
   - Expand documentation with detailed guides for setup, advanced examples, and best practices.

5. **Testing**:
   - Add unit and integration tests to ensure connector robustness in various scenarios.

---
