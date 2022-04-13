package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	store "yandex-praktikum/internal/store"
)

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

	responseURL := store.GetShortURL(urlToShort, r.Host)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(responseURL))
	fmt.Fprint(w)
}

func PostShorten(w http.ResponseWriter, r *http.Request) {
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

	type in struct {
		URL string `json:"url,omitempty"`
	}

	valueIn := in{}

	if err := json.Unmarshal([]byte(body), &valueIn); err != nil {
		http.Error(w, "Shorten unmarshal error", http.StatusBadRequest)
		return
	}

	if valueIn.URL == "" {
		http.Error(w, "Shorten url not found", http.StatusBadRequest)
		return
	}

	responseURL := store.GetShortURL(valueIn.URL, r.Host)

	type out struct {
		Result string `json:"result"`
	}

	valueOut := out{
		Result: responseURL,
	}

	result, err := json.Marshal(valueOut)
	if err != nil {
		http.Error(w, "Shorten marshal error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
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

	reternURL, err := store.GetURL(id)
	if err != "" {
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Add("Location", reternURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write([]byte(""))
	fmt.Fprint(w)

}
