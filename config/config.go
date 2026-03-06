package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config holds the application configuration.
// Values are populated from environment variables.
type Config struct {
	ZmqBindAddress string `envconfig:"ZMQ_BIND_ADDRESS" default:"tcp://127.0.0.1:5555"`
	ZmqPubAddress  string `envconfig:"ZMQ_PUB_ADDRESS" default:"tcp://127.0.0.1:5556"`
}

// Load reads configuration from environment variables and returns a new Config struct.
func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
