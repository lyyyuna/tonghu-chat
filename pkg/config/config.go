package config

import (
	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
)

type Config struct {
	Nats     Nats
	Redis    Redis
	ChatPort int
	ChatHost string
}

type Nats struct {
	ClusterId string
	ClientId  string
	Host      string
	Port      int
}

type Redis struct {
	Host string
	Port int
}

func ReadConfig(configPath string) *Config {
	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		zap.S().Fatalf("cannot read the config, the err: %v", err)
	}

	return &config
}
