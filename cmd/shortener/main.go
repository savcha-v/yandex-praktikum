package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

var Urls map[string]string = map[string]string{
	"0": "",
}
var Unic int

func PostShort(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlToShort := string(body)
	if urlToShort == "" {
		http.Error(w, "Shortcut url not found", http.StatusBadRequest)
		return
	}

	Unic += 1
	UnicStr := strconv.Itoa(Unic)
	Urls[UnicStr] = urlToShort

	responseURL := "http://" + r.Host + r.URL.String() + "?id=" + UnicStr

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(responseURL))
	fmt.Fprint(w)
}

func GetShort(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "'id' missing", http.StatusBadRequest)
		return
	}

	reternURL, exists := Urls[id]
	if !exists {
		http.Error(w, "'id' not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Add("Location", reternURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write([]byte(""))
	fmt.Fprint(w)

}

func main() {

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/", GetShort)
		r.Post("/", PostShort)
	})

	delete(Urls, "0")

	log.Fatal(http.ListenAndServe(":8080", r))

}
