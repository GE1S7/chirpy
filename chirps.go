package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GE1S7/m/v2/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	r.PathValue
	cfg.dbQueries.

func (cfg *apiConfig) getAllChirpsHandler(w http.ResponseWriter, r *http.Request) {

	dbChirps, err := cfg.dbQueries.GetChirpsCreatedAsc(r.Context())
	if err != nil {
		respondWithError(w, 500, "Error fetching resource")
	} else {
		var chirps []Chirp
		for _, e := range dbChirps {

			chirp := Chirp{
				ID:        e.ID,
				CreatedAt: e.CreatedAt,
				UpdatedAt: e.UpdatedAt,
				Body:      e.Body,
				UserID:    e.UserID,
			}

			chirps = append(chirps, chirp)

		}
		respondWithJson(w, 200, chirps)
	}

}

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	var parameters Parameters

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&parameters); err != nil {
		respondWithError(w, 400, "Error processing a request")
	}

	cleanBody, err := validateChirp(w, r, parameters.Body)

	if err != nil {
		respondWithError(w, 400, "Error validating chirp")
		fmt.Println("Error validating chirp:", err)
	} else {
		chirpParams := database.CreateChirpParams{
			Body:   cleanBody,
			UserID: parameters.UserID,
		}

		createdChirp, err := cfg.dbQueries.CreateChirp(r.Context(), chirpParams)
		if err != nil {
			fmt.Println("Error creating chirp", err)
			respondWithError(w, 500, "Error creating chirp")
		} else {
			chirp := Chirp{
				ID:        createdChirp.ID,
				CreatedAt: createdChirp.CreatedAt,
				UpdatedAt: createdChirp.UpdatedAt,
				Body:      createdChirp.Body,
				UserID:    createdChirp.UserID,
			}

			respondWithJson(w, 201, chirp)
		}
	}
}
