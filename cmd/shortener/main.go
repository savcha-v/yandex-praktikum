package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	handlers "yandex-praktikum/internal/handlers"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func createServer(port int, cfg Config) *http.Server {

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

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		server := createServer(8080, cfg)
		log.Fatal(server.ListenAndServe())
		wg.Done()
	}()

	go func() {
		port, err := strconv.Atoi(cfg.ServerAddress)
		if err != nil {
			log.Fatal()
		}
		server := createServer(port, cfg)
		log.Fatal(server.ListenAndServe())
		wg.Done()
	}()

	wg.Wait()

}
