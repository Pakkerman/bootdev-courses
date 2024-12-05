package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pakkerman/bootdev-chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
}

func main() {
	godotenv.Load(".env")
	dbURL := os.Getenv("DB_URL")

	fmt.Println("DB URL", dbURL)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Errorf("something wrong when connecting to database: %v", err)
		return
	}

	const fsRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      database.New(db),
		platform:       os.Getenv("PLATFORM"),
	}

	// This will redirect .../app/ to ./index.html
	// I think it takes the url, lcoalhost:8080/app/ and strip out the app part,
	// and serve the second part which in thic case is ".",
	// and by default, the index.html file
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(fsRoot)))))
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handleMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handleReset)
	mux.HandleFunc("POST /api/users", cfg.handleUsers)
	mux.HandleFunc("POST /api/chirps", cfg.handleChirps)

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("listening on port 8080")
	log.Fatal(s.ListenAndServe())
}
