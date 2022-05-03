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
	r.Use(cookie.SetUserID)
	r.Route("/", func(r chi.Router) {

		r.Get("/"+cfg.BaseURL+"/", handlers.GetShort)
		r.Post("/", handlers.PostShort(cfg))
		r.Post("/api/shorten", handlers.PostShorten(cfg))
		r.Get("/api/user/urls", handlers.GetUserShorts)

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
