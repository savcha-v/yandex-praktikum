package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

var typeToCompress string

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func init() {
	typeToCompress =
		`application/javascript
		application/json
		text/css
		text/html
		text/plain
		application/gzip
		application/x-gzip
		text/xml`
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func CompressHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

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
		if strings.Contains(typeToCompress, r.Header.Get("Content-Type")) {

			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				io.WriteString(w, err.Error())
			}
			defer gz.Close()
			w.Header().Set("Content-Encoding", "gzip")

			next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)

		}
	})
}
