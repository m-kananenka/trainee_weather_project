package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Token string `envconfig:"token"`

	Port string `envconfig:"app_port"`
}

func (c *Config) Process() error {
	return envconfig.Process("example", c)
}
