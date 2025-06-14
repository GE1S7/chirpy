package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/GE1S7/chirpy/internal/auth"
	"github.com/GE1S7/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpID")
	idParsed, err := uuid.Parse(chirpId)

	if chirpId == "" || err != nil {
		respondWithError(w, 404, "error parsing path")
	} else {
		dbChirp, err := cfg.dbQueries.GetChirp(r.Context(), idParsed)
		if err != nil {
			respondWithError(w, 404, "error fetching resource")
		}

		chirp := Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}

		respondWithJson(w, 200, chirp)

	}
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	var dbChirps []database.Chirp
	var err error
	userID, err := uuid.Parse(r.URL.Query().Get("author_id"))
	if err != nil {
		dbChirps, err = cfg.dbQueries.GetChirpsCreatedAsc(r.Context())
		if err != nil {
			respondWithError(w, 500, "Error fetching resource")
			return
		}
	} else {
		dbChirps, err = cfg.dbQueries.GetChirpsByAuthor(r.Context(), userID)
		if err != nil {
			respondWithError(w, 500, "Error fetching resource")
			return
		}
	}

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

	if r.URL.Query().Get("sort") == "desc" {
		slices.Reverse(chirps)
	}

	respondWithJson(w, 200, chirps)

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

	// validate jwt
	var userid uuid.UUID
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 400, "Error processing a request")
	} else {
		userid, err = auth.ValidateJWT(token, cfg.jwtSecret)
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println("token:", token)
			respondWithError(w, 401, "Unauthorized")
		} else {

			cleanBody, err := validateChirp(w, r, parameters.Body)

			if err != nil {
				respondWithError(w, 400, "Error validating chirp")
				fmt.Println("Error validating chirp:", err)
			} else {
				chirpParams := database.CreateChirpParams{
					Body:   cleanBody,
					UserID: userid,
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
	}
}

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 401, "Error parsing path")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Error extracting token from header")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Error validating token")
		return
	}

	chirpData, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "Error fetching chirp from the database")
		return
	}

	if chirpData.UserID != userID {
		respondWithError(w, 403, "UserID in request does not match chirp.UserID")
		return
	}

	err = cfg.dbQueries.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 401, "Error deleting chirp")
		return
	}

	respondWithJson(w, 204, "")

}
