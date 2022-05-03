package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"yandex-praktikum/internal/config"
	"yandex-praktikum/internal/cookie"
	store "yandex-praktikum/internal/store"
)

func PostShort(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

		userID := cookie.GetUserID(r)

		responseURL := store.GetShortURL(urlToShort, r.Host, cfg, userID)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(responseURL))
		fmt.Fprint(w)
	}
}

func PostShorten(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

		userID := cookie.GetUserID(r)

		responseURL := store.GetShortURL(valueIn.URL, r.Host, cfg, userID)

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
}

func GetShort(w http.ResponseWriter, r *http.Request) {

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

func GetUserShorts(w http.ResponseWriter, r *http.Request) {

	userID := cookie.GetUserID(r)

	userShorts := store.GetUserShorts(userID)

	w.Header().Set("Content-Type", "application/json")
	if len(userShorts) != 0 {
		result, err := json.Marshal(userShorts)
		if err != nil {
			http.Error(w, "Shorten marshal error", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}

	fmt.Fprint(w)
}
