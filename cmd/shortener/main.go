package main

import (
	"flag"
	"log"
	"net/http"
	config "yandex-praktikum/internal/config"
	handlers "yandex-praktikum/internal/handlers"
	"yandex-praktikum/internal/store"

	"github.com/go-chi/chi/v5"
)

func createServer() *http.Server {

	flag.Parse()

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/"+config.Cfg.BaseURL+"/", handlers.GetShort)
		r.Post("/", handlers.PostShort)
		r.Post("/api/shorten", handlers.PostShorten)
	})

	server := http.Server{
		Addr:    config.Cfg.ServerAddress,
		Handler: r,
	}

	return &server
}

func main() {

	server := createServer()
	store.InitStorage()
	log.Fatal(server.ListenAndServe())

}
