package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {

	type Params struct {
		Email string `json:"email"`
	}

	var params Params

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {

		respondWithError(w, 400, "Error decoding request")

	} else {

		record, err := cfg.dbQueries.CreateUser(r.Context(), params.Email)
		if err != nil {
			respondWithError(w, 400, "Error while fetching from database")
		}

		user := User{
			ID:        record.ID,
			CreatedAt: record.CreatedAt,
			UpdatedAt: record.UpdatedAt,
			Email:     record.Email,
		}
		respondWithJson(w, 201, user)
	}

}
