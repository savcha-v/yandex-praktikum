package main

import (
	"log"
	"net/http"
	handlers "yandex-praktikum/internal/handlers"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"BASE_URL"`
}

func createServer() *http.Server {

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	//port, _ := strconv.Atoi(cfg.ServerAddress)

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/"+cfg.BaseURL+"/", handlers.GetShort)
		r.Post("/", handlers.PostShort)
		r.Post("/api/shorten", handlers.PostShorten)
	})

	server := http.Server{
		Addr:    cfg.ServerAddress,
		Handler: r,
	}

	return &server
}

func main() {

	server := createServer()
	log.Fatal(server.ListenAndServe())

}
