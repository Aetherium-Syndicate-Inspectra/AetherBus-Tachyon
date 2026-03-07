package config

import "os"

// Config holds the application configuration.
// Values are populated from environment variables.
type Config struct {
	ZmqBindAddress string
	ZmqPubAddress  string
}

// Load reads configuration from environment variables and returns a new Config struct.
func Load() (*Config, error) {
	cfg := &Config{
		ZmqBindAddress: getenvOrDefault("ZMQ_BIND_ADDRESS", "tcp://127.0.0.1:5555"),
		ZmqPubAddress:  getenvOrDefault("ZMQ_PUB_ADDRESS", "tcp://127.0.0.1:5556"),
	}

	return cfg, nil
}

func getenvOrDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}

	return defaultValue
}
