package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"BASE_URL"`
	FileStor      string `env:"FILE_STORAGE_PATH"`
}

var cfg Config

func init() {

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "")
	flag.StringVar(&cfg.FileStor, "f", cfg.FileStor, "")
}

func BaseURL() string {
	return cfg.BaseURL
}
func ServerAddress() string {
	return cfg.ServerAddress
}
func FileStor() string {
	return cfg.FileStor
}
