package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const fsRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// This will redirect .../app/ to ./index.html
	// I think it takes the url, lcoalhost:8080/app/ and strip out the app part,
	// and serve the second part which in thic case is ".",
	// and by default, the index.html file
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(fsRoot)))))
	mux.HandleFunc("GET /healthz", handleReadiness)
	mux.HandleFunc("GET /metrics", cfg.handleMetrics)
	mux.HandleFunc("POST /reset", cfg.handleReset)

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("listening on port 8080")
	log.Fatal(s.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}
