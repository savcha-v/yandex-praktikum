package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

var Urls map[string]string = map[string]string{
	"0": "",
}
var Unic int

// Short — обработчик запроса.
func Short(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:

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

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(UnicStr))
		fmt.Fprintln(w)

	case http.MethodGet:

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "'id' missing", http.StatusBadRequest)
			return
		}

		reternUrl, exists := Urls[id]
		if !exists {
			http.Error(w, "'id' not found", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Add("Location", reternUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte(""))
		fmt.Fprintln(w)

	default:

		http.Error(w, "POST or GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
}

func main() {

	http.HandleFunc("/", Short)

	delete(Urls, "0")

	log.Fatal(http.ListenAndServe(":8080", nil))

}
