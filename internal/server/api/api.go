package api

import (
	"compress/flate"
	"compress/gzip"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/pochtalexa/go-cti-middleware/internal/server/handlers"
	"github.com/rs/zerolog/log"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// проверяем, что клиент отправил серверу сжатые данные в формате gzip
func checkGzipEncoding(r *http.Request) bool {

	encodingSlice := r.Header.Values("Content-Encoding")
	encodingsStr := strings.Join(encodingSlice, ",")
	encodings := strings.Split(encodingsStr, ",")

	log.Info().Str("encodingsStr", encodingsStr).Msg("checkGzipEncoding")

	for _, el := range encodings {
		if el == "gzip" {
			return true
		}
	}

	return false
}

func GzipDecompression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if checkGzipEncoding(r) {
			gzipReader, err := gzip.NewReader(r.Body)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			r.Body = gzipReader
			defer gzipReader.Close()
		}

		log.Info().Msg("GzipDecompression passed")

		next.ServeHTTP(w, r)
	})
}

func RunAPI(urlStr string) error {
	logger := httplog.NewLogger("httplog", httplog.Options{
		LogLevel: slog.LevelDebug,
		//JSON:             true,
		Concise:          false,
		RequestHeaders:   true,
		ResponseHeaders:  true,
		MessageFieldName: "msg",
		//LevelFieldName:   "severity",
		TimeFieldFormat: time.RFC3339,
		Tags: map[string]string{
			"version": "v1.0",
			"env":     "dev",
		},
		//QuietDownRoutes: []string{
		//	"/",
		//	"/ping",
		//},
		//QuietDownPeriod: 10 * time.Second,
		//SourceFieldName: "source",
	})

	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.Recoverer)
	mux.Use(GzipDecompression)
	mux.Use(middleware.Compress(flate.DefaultCompression, "application/json", "text/html"))
	mux.Use(middleware.Logger)
	mux.Use(httplog.RequestLogger(logger))

	mux.Get("/", handlers.RootHandler)

	log.Info().Str("Running on", urlStr).Msg("http server started")

	return http.ListenAndServe(urlStr, mux)
}
