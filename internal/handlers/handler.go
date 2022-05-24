package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

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

		userID := cookie.GetUserID(r, cfg)

		responseURL, httpStatus := store.GetShortURL(r.Context(), urlToShort, r.Host, cfg, userID)

		fmt.Fprintln(os.Stdout, "PostShort: "+responseURL)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(httpStatus)
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

		if err := json.Unmarshal(body, &valueIn); err != nil {
			http.Error(w, "Shorten unmarshal error", http.StatusBadRequest)
			return
		}

		if valueIn.URL == "" {
			http.Error(w, "Shorten url not found", http.StatusBadRequest)
			return
		}

		userID := cookie.GetUserID(r, cfg)

		responseURL, httpStatus := store.GetShortURL(r.Context(), valueIn.URL, r.Host, cfg, userID)

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
		w.WriteHeader(httpStatus)
		w.Write(result)
		fmt.Fprint(w)
	}
}

func GetShort(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "'id' missing", http.StatusBadRequest)
			return
		}

		reternURL, err := store.GetURL(r.Context(), id, cfg)
		if err != "" {
			http.Error(w, err, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		if reternURL != "" {
			w.Header().Add("Location", reternURL)
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			w.WriteHeader(http.StatusGone)
		}

		w.Write([]byte(""))
		fmt.Fprint(w)

	}
}

func GetUserShorts(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := cookie.GetUserID(r, cfg)

		userShorts := store.GetUserShorts(r.Context(), cfg, userID)

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
}

func GetPing(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		statusPing := store.PingDB(r.Context(), cfg.ConnectDB)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(statusPing)
		w.Write([]byte(""))
		fmt.Fprint(w)
	}
}

func PostBatch(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var v []store.RequestURL
		userID := cookie.GetUserID(r, cfg)

		if err := json.Unmarshal(body, &v); err != nil {
			http.Error(w, "Batch unmarshal error", http.StatusBadRequest)
			return
		}

		valueOut := store.ShortURLs(r.Context(), v, r.Host, cfg, userID)

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

func DeleteURLs(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Fprintln(os.Stdout, "DeleteURLs body: "+string(body[:]))

		var v []string

		if err := json.Unmarshal(body, &v); err != nil {
			http.Error(w, "delete urls unmarshal error", http.StatusBadRequest)
			return
		}

		for _, un := range v {
			fmt.Fprintln(os.Stdout, un)
		}

		userID := cookie.GetUserID(r, cfg)

		strDel := config.StructToDelete{
			UserID: userID,
			ListID: v,
		}

		go func() {
			cfg.DeleteChan <- strDel
		}()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(""))
		fmt.Fprint(w)

	}
}
