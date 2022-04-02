package handlerfuncs

import (
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

	id := store.GetID(urlToShort)

	responseURL := "http://" + r.Host + r.URL.String() + "?id=" + id

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
