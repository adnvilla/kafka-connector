package kafka_connector

type Provider int

const (
	ZKakfa Provider = iota
)

type Config struct {
	// BootstrapServers is a list of broker addresses
	BootstrapServers []string
	ClientID         string
	Provider         Provider
	UseGlobalClient  bool
}
