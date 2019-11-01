package config

import (
	"os"
)

type Config struct {
	AmqpUri           string `toml:"amqpUri"`
	ConsumerPoolSize  int    `toml:"consumerPoolSize"`
	PublisherPoolSize int    `toml:"publisherPoolSize"`
	IdleTimeout       int    `toml:"idleTimeout"` // Pool idle timeout in seconds
	HttpTimeout       int    `toml:"httpTimeout"`
	HttpPoolSize      int    `toml:"httpPoolSize"`
	RedisUri          string `toml:"redisUri"`
	RedisPoolSize     int    `toml:"redisPoolSize"`
}

type Provider int

const (
	Toml    Provider = iota
	EnvVars Provider = iota
)

var Conf Config

func (p *Provider) LoadConfig() {
	switch *p {
	case Toml:
		loadTomlConfig()
	case EnvVars:
		loadEnvVarConfig()
	}
}

func GetConfigProvider() Provider {
	tomlExists := false
	if _, err := os.Stat("config.toml"); err == nil {
		tomlExists = true
	}

	if tomlExists {
		return Toml
	} else {
		return EnvVars
	}
}
