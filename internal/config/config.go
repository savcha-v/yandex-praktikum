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

var Cfg Config

var ServerAddress string

func init() {

	err := env.Parse(&Cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := env.Parse(&Cfg); err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&Cfg.ServerAddress, "a", Cfg.ServerAddress, "")
	flag.StringVar(&Cfg.BaseURL, "b", Cfg.BaseURL, "")
	flag.StringVar(&Cfg.FileStor, "f", Cfg.FileStor, "")
}
