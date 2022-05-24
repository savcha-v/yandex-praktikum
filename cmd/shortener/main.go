package main

import (
	"log"
	"net/http"
	"yandex-praktikum/internal/compress"
	config "yandex-praktikum/internal/config"
	"yandex-praktikum/internal/cookie"
	handlers "yandex-praktikum/internal/handlers"
	"yandex-praktikum/internal/store"

	"github.com/go-chi/chi/v5"
)

func createServer(cfg config.Config) *http.Server {

	r := chi.NewRouter()
	r.Use(compress.CompressHandler)
	r.Use(cookie.SetUserID(cfg))

	r.Route("/", func(r chi.Router) {

		r.Post("/", handlers.PostShort(cfg))
		r.Post("/api/shorten", handlers.PostShorten(cfg))
		r.Post("/api/shorten/batch", handlers.PostBatch(cfg))
		r.Get("/api/user/urls", handlers.GetUserShorts(cfg))
		r.Get("/ping", handlers.GetPing(cfg))
		r.Get("/"+cfg.BaseURL+"/", handlers.GetShort(cfg))
		r.Delete("/api/user/urls", handlers.DeleteURLs(cfg))
	})

	server := http.Server{
		Addr:    cfg.ServerAddress,
		Handler: r,
	}

	return &server
}

func main() {

	cfg := config.NewConfig()
	store.InitStorage(&cfg)
	defer cfg.ConnectDB.Close()
	defer close(cfg.DeleteChan)
	server := createServer(cfg)
	log.Fatal(server.ListenAndServe())

}
