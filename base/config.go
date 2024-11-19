package base

type Config struct {
	// BootstrapServers is a list of broker addresses
	BootstrapServers []string
	ClientID         string
	UseGlobalClient  bool
}
