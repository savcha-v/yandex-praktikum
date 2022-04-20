package main

import (
	"log"
	"net/http"
	"yandex-praktikum/internal/compress"
	config "yandex-praktikum/internal/config"
	handlers "yandex-praktikum/internal/handlers"
	"yandex-praktikum/internal/store"

	"github.com/go-chi/chi/v5"
)

func createServer(cfg config.Config) *http.Server {

	r := chi.NewRouter()
	r.Use(compress.CompressHandler)
	r.Route("/", func(r chi.Router) {

		r.Get("/"+cfg.BaseURL+"/", handlers.GetShort)
		r.Post("/", handlers.PostShort(cfg))
		r.Post("/api/shorten", handlers.PostShorten(cfg))

	})

	server := http.Server{
		Addr:    cfg.ServerAddress,
		Handler: r,
	}

	return &server
}

func main() {

	cfg := config.NewConfig()
	server := createServer(cfg)
	store.InitStorage(cfg.FileStor)

	log.Fatal(server.ListenAndServe())

}
