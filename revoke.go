package main

import (
	"net/http"

	"github.com/GE1S7/chirpy/internal/auth"
)

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "No token in request header found")
		return
	}

	err = cfg.dbQueries.RevokeToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "Error revoking token")
		return
	}

	respondWithJson(w, 204, "")

}
