package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/GE1S7/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "No refresh token found in request header")
		return
	}

	tokenData, err := cfg.dbQueries.GetTokenData(r.Context(), refreshToken)

	if tokenData.RevokedAt.Valid == true {
		respondWithError(w, 401, fmt.Sprintf("%v", err))
		return
	}

	expiresAt := tokenData.ExpiresAt
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("%v", err))
		return
	}

	if expiresAt.Before(time.Now()) {
		respondWithError(w, 401, "Expired refresh token")
		fmt.Println("Expired!")
		return
	}

	accessExpiresIn, err := time.ParseDuration("1h")
	if err != nil {
		respondWithError(w, 401, "")
	}

	accessToken, err := auth.MakeJWT(tokenData.UserID, cfg.jwtSecret, accessExpiresIn)
	if err != nil {
		respondWithError(w, 401, "Error creating access token")
	}

	type Resp struct {
		Token string `json:"token"`
	}

	resp := Resp{
		Token: accessToken,
	}

	respondWithJson(w, 200, resp)

}
