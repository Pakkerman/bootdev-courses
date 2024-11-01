package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) handleUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json: "email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(500)

		res := struct {
			Error string `json:"error"`
		}{
			Error: fmt.Sprintf("Something wrong with server: %v", err),
		}

		data, err := json.Marshal(res)
		if err != nil {
			log.Printf("Error with marshalling JSON: %s", err)
		}
		w.Write(data)
		return
	}

	email := params.Email
	user, err := cfg.dbQueries.CreateUser(r.Context(), email)
	if err != nil {
		{
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(500)

			res := struct {
				Error string `json:"error"`
			}{
				Error: fmt.Sprintf("Something wrong with createUser with DB: %v", err),
			}

			data, err := json.Marshal(res)
			if err != nil {
				log.Printf("Error with marshalling JSON: %s", err)
			}
			w.Write(data)
			return
		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(201)
	res := struct {
		Id        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Email     string `json:"email"`
	}{
		Id:        user.ID.UUID.String(),
		CreatedAt: user.CreatedAt.Time.String(),
		UpdatedAt: user.UpdatedAt.Time.String(),
		Email:     user.Email,
	}

	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error with marshalling JSON: %s", err)
	}
	w.Write(data)
}
