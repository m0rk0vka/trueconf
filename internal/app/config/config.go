package config

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

const defaultServerPort = 8080

// Config represents an application configuration.
type Config struct {
	// the server port. Defaults to 8080
	ServerPort int `yaml:"server_port" env:"SERVER_PORT"`
}

func Load(file string /*logger log.Logger*/) (*Config, error) {
	// default config
	c := Config{
		ServerPort: defaultServerPort,
	}

	// load from YAML config file
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}

	// TODO load from environment variables prefixed with "APP_"

	return &c, nil
}
