package config

import (
	"database/sql"
	"flag"
	"log"

	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"BASE_URL"`
	FileStor      string `env:"FILE_STORAGE_PATH"`
	DataBase      string `env:"DATABASE_DSN"`
	ConnectDB     *sql.DB
	Key           string
	DeleteChan    chan StructToDelete
}

type StructToDelete struct {
	UserID string
	ListID []string
}

func NewConfig() Config {

	var cfg Config

	cfg.Key = "10c57de0"

	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "")
	flag.StringVar(&cfg.FileStor, "f", cfg.FileStor, "")
	flag.StringVar(&cfg.DataBase, "d", cfg.DataBase, "")

	flag.Parse()

	cfg.DeleteChan = make(chan StructToDelete, 10)

	return cfg
}
