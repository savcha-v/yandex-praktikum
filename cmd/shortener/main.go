package main

import (
	"flag"
	"log"
	"net/http"
	"yandex-praktikum/internal/compress"
	config "yandex-praktikum/internal/config"
	handlers "yandex-praktikum/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func createServer() *http.Server {

	flag.Parse()

	r := chi.NewRouter()
	r.Use(compress.CompressHandler)
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
