package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	mux.Handle("/", http.FileServer(http.Dir(".")))
	fmt.Println("listening on port 8080")
	s.ListenAndServe()
}
