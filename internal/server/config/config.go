package config

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog/log"
	"os"
)

type Config struct {
	Settings Settings
	CtiAPI   CtiAPI
	HttpAPI  HttpAPI
}

type Settings struct {
	LogLevel string
}

type CtiAPI struct {
	Scheme string
	Path   string
	Host   string
	Port   string
}

type HttpAPI struct {
	Scheme string
	Host   string
	Port   string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) ReadConfigFile() error {
	fileName := "config.toml"

	file, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("ReadFile: %w", err)
	}

	if err := toml.Unmarshal(file, c); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	log.Info().Msg("config file parsed - ok")

	return nil
}
