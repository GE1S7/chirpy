package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/GE1S7/chirpy/internal/auth"
	"github.com/GE1S7/chirpy/internal/database"
	"github.com/google/uuid"
)

type UserOut struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {

	type Params struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	var params Params

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {

		respondWithError(w, 400, "Error decoding request")

	} else {

		hashed, err := auth.HashPassword(params.Password)
		if err != nil {
			respondWithError(w, 401, "")
		}

		userParams := database.CreateUserParams{
			Email: params.Email,
			HashedPassword: sql.NullString{
				String: hashed,
				Valid:  true,
			},
		}

		record, err := cfg.dbQueries.CreateUser(r.Context(), userParams)
		if err != nil {
			respondWithError(w, 400, "Error while fetching from database")
		}

		user := UserOut{
			ID:        record.ID,
			CreatedAt: record.CreatedAt,
			UpdatedAt: record.UpdatedAt,
			Email:     record.Email,
		}
		respondWithJson(w, 201, user)
	}

}
