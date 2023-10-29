package config

import "github.com/kelseyhightower/envconfig"

// MustLoadConfig Load configuration from environment variables.
func MustLoadConfig() Config {
	var c Config
	envconfig.MustProcess("AUTH", &c)

	return c
}
