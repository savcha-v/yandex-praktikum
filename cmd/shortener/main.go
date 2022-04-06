package main

import (
	"log"
	"net/http"
	handlers "yandex-praktikum/internal/handlers"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func main() {

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/"+cfg.BaseURL+"/", handlers.GetShort)
		r.Post("/", handlers.PostShort)
		r.Post("/api/shorten", handlers.PostShorten)
	})

	log.Fatal(http.ListenAndServe(cfg.ServerAddress+":8080", r))

}
