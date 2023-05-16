package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	APIKey string `envconfig:"api_key"`
	URL    string `envconfig:"url"`
}

func (c *Config) Process() error {
	return envconfig.Process("example", c)
}
