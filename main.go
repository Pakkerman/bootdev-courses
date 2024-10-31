package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handleMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handleReset)
	mux.HandleFunc("POST /api/validate_chirp", handleValidate)

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("listening on port 8080")
	log.Fatal(s.ListenAndServe())
}

func handleValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type responseError struct {
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(500)

		res := responseError{
			Error: "Something went wrong",
		}

		data, err := json.Marshal(res)
		if err != nil {
			log.Printf("Error with marshalling JSON: %s", err)
		}

		w.Write(data)
		return
	}

	content := params.Body

	if len(content) > 140 {
		{
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(400)

			res := responseError{
				Error: "Chirp is too long",
			}

			data, err := json.Marshal(res)
			if err != nil {
				log.Printf("Error with marshalling JSON: %s", err)
			}

			w.Write(data)
			return
		}
	}

	type responseBody struct {
		Body string `json:"cleaned_body"`
	}

	w.Header().Add("Content-Type", "application/json")
	res := responseBody{
		Body: chirpFilter(content),
	}

	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error with marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Write(data)
}

func chirpFilter(s string) string {
	cenceredWords := []string{
		"profane", "kerfuffle", "sharbert", "fornax",
	}

	arr := strings.Split(strings.ToLower(s), " ")

	for i := 0; i < len(arr); i++ {
		for k := 0; k < len(cenceredWords); k++ {
			if arr[i] == cenceredWords[k] {
				arr[i] = "****"
			}
		}
	}

	fmt.Println(arr)

	return "1"
}
