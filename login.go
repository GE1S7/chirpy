package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GE1S7/chirpy/internal/auth"
	"github.com/GE1S7/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type UserOut struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}

	var params Params

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {

		respondWithError(w, 400, "")
		return

	} else {

		user, err := cfg.dbQueries.GetUserByMail(r.Context(), params.Email)

		if err != nil {

			respondWithError(w, 401, fmt.Sprintf("%v", err))
			return

		} else {
			err := auth.CheckPasswordHash(user.HashedPassword.String, params.Password)
			if err != nil {

				respondWithError(w, 401, fmt.Sprintf("%v", err))
			} else {
				var expiresIn time.Duration

				expiresIn, err = time.ParseDuration("1h")
				if err != nil {
					respondWithError(w, 401, fmt.Sprintf("%v", err))
					return
				}

				// create access token
				token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiresIn)
				fmt.Println("created token:", token)
				if err != nil {
					respondWithError(w, 401, fmt.Sprintf("%v", err))
					return
				}

				// create refresh token
				refreshToken, err := auth.MakeRefreshToken()
				if err != nil {
					respondWithError(w, 401, fmt.Sprintf("%v", err))
					return
				}

				expiresInRefresh, err := time.ParseDuration(fmt.Sprintf("%vh", 60*24))
				if err != nil {
					respondWithError(w, 401, fmt.Sprintf("%v", err))
					return
				}

				refreshTokenParams := database.CreateRefreshTokenParams{
					Token:     refreshToken,
					UserID:    user.ID,
					ExpiresAt: time.Now().UTC().Add(expiresInRefresh),
				}

				cfg.dbQueries.CreateRefreshToken(r.Context(), refreshTokenParams)

				userOut := UserOut{
					ID:           user.ID,
					CreatedAt:    user.CreatedAt,
					UpdatedAt:    user.UpdatedAt,
					Email:        user.Email,
					Token:        token,
					RefreshToken: refreshToken,
				}
				respondWithJson(w, 200, userOut)
			}

		}

	}

}
