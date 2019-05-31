package config

import (
	"os"
)

type Config struct {
	AmqpUri     string `toml:"amqpUri"`
	PoolSize    int    `toml:"poolSize"`   // Amount of connections to be BLPOP'ing
	IdleTimeout int    `toml:"idleTimout"` // Pool idle timeout in seconds
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
