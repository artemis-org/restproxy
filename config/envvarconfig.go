package config

import (
	"github.com/kelseyhightower/envconfig"
)

func loadEnvVarConfig() {
	if err := envconfig.Process("artemis-org", &Conf); err != nil {
		panic(err)
	}
}
