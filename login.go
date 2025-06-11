package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GE1S7/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
	}

	type UserOut struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	var params Params

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {

		respondWithError(w, 400, "")

	} else {

		user, err := cfg.dbQueries.GetUserByMail(r.Context(), params.Email)

		if err != nil {

			respondWithError(w, 401, "")

		} else {
			err := auth.CheckPasswordHash(user.HashedPassword.String, params.Password)
			if err != nil {

				respondWithError(w, 401, "")
				fmt.Println(user)
				fmt.Println(params.Password)
			} else {
				userOut := UserOut{
					ID:        user.ID,
					CreatedAt: user.CreatedAt,
					UpdatedAt: user.UpdatedAt,
					Email:     user.Email,
				}
				respondWithJson(w, 200, userOut)
			}

		}

	}

}
