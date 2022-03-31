package main

import (
	"log"
	"net/http"
	"yandex-praktikum/cmd/shortener/handlerfuncs"

	"github.com/go-chi/chi/v5"
)

func main() {

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlerfuncs.GetShort)
		r.Post("/", handlerfuncs.PostShort)
	})

	delete(handlerfuncs.Urls, "0")

	log.Fatal(http.ListenAndServe(":8080", r))

}
