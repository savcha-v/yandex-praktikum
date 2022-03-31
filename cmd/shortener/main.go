package main

import (
	"log"
	"net/http"
	"yandex-praktikum/cmd/shortener/handlerFuncs"

	"github.com/go-chi/chi/v5"
)

func main() {

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlerFuncs.GetShort)
		r.Post("/", handlerFuncs.PostShort)
	})

	delete(handlerFuncs.Urls, "0")

	log.Fatal(http.ListenAndServe(":8080", r))

}
