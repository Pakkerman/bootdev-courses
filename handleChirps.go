package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handleChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string `json:"body"`
		UserId string `json:"user_id"`
	}

	dec := json.NewDecoder(r.Body)
	params := parameters{}
	err := dec.Decode(&params)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("something wrong with decoding body"))
	}

	content := params.Body
	userId := params.UserId
	if len(params.Body) > 140 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(400)

		res := struct {
			Error string `json:"error"`
		}{
			Error: "Chirp is too long",
		}

		data, err := json.Marshal(res)
		if err != nil {
			log.Printf("Error with marshalling JSON: %s", err)
		}

		w.Write(data)
		return
	}

	content = chirpFilter(content)
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
