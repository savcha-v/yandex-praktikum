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

func NewConfig() Config {

	var cfg Config

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

	flag.Parse()

	return cfg
}
