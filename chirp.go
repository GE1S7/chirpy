package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GE1S7/m/v2/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) chirpHandler(w http.ResponseWriter, r *http.Request) {
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
			Body: cleanBody,
			UserID: uuid.NullUUID{
				UUID:  parameters.UserID,
				Valid: true,
			},
		}

		createdChirp, err := cfg.dbQueries.CreateChirp(r.Context(), chirpParams)
		if err != nil {
			fmt.Println("Error creating chirp")
			respondWithError(w, 500, "Error creating chirp")
		} else {
			respondWithJson(w, 201, createdChirp)
		}

	}

}
