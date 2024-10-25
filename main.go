package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const fsRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	// This will redirect .../app/ to ./index.html
	// I think it takes the url, lcoalhost:8080/app/ and strip out the app part,
	// and serve the second part which in thic case is ".",
	// and by default, the index.html file
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(fsRoot))))
	mux.HandleFunc("/healthz", handleReadiness)

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("listening on port 8080")
	log.Fatal(s.ListenAndServe())
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
