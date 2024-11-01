package main

import (
	"net/http"
)

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(403)
		w.Write([]byte("403 Forbidden"))
	}

	// Reset traffic counter
	cfg.fileserverHits.Store(0)

	// Reset user table
	_, err := cfg.dbQueries.TruncateUserTable(r.Context())
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("something wrong with server when truncating user table"))
	}

	w.WriteHeader(200)
	w.Write([]byte("reset"))
}
