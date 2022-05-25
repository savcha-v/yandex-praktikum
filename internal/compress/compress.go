package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

const typeToCompress = `application/javascript
		application/json
		text/css
		text/html
		text/plain
		application/gzip
		application/x-gzip
		text/xml`

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func CompressHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// распаковать
		if r.Method == http.MethodPost && r.Header.Get(`Content-Encoding`) == "gzip" {

			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			r.Body = gz
			defer gz.Close()

		}

		// запаковать
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") &&
			strings.Contains(typeToCompress, r.Header.Get("Content-Type")) {

			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				io.WriteString(w, err.Error())
			}
			defer gz.Close()
			w.Header().Set("Content-Encoding", "gzip")

			w = gzipWriter{ResponseWriter: w, Writer: gz}

		}

		next.ServeHTTP(w, r)

	})
}
