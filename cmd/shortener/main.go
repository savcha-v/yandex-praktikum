package main

import (
	"fmt"
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

func createServer() *http.Server {

	port := 8080
	var cfg Config
	env.Parse(&cfg)
	//port, _ = strconv.Atoi(cfg.ServerAddress)

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/"+cfg.BaseURL+"/", handlers.GetShort)
		r.Post("/", handlers.PostShort)
		r.Post("/api/shorten", handlers.PostShorten)
	})

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: r,
	}

	return &server
}

func main() {

	server := createServer()
	log.Fatal(server.ListenAndServe())

}
