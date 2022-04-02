package main

import (
	"log"
	"net/http"
	handlers "yandex-praktikum/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func main() {

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.GetShort)
		r.Post("/", handlers.PostShort)
	})

	log.Fatal(http.ListenAndServe(":8080", r))

}
