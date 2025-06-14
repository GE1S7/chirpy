package main

import (
	"fmt"
	"io"
	"net/http"

	"encoding/json"

	"github.com/GE1S7/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) upgradeUserHandler(w http.ResponseWriter, r *http.Request) {
	type ReqData struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("%v", err))
		return
	}

	if key != cfg.polkaApiKey {
		respondWithError(w, 401, fmt.Sprintf("Incorrect api key"))
		return
	}

	reqData := ReqData{}

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("%v", err))
		return
	}

	if err := json.Unmarshal(reqBytes, &reqData); err != nil {
		respondWithError(w, 401, fmt.Sprintf("%v", err))
		return
	}

	if reqData.Event != "user.upgraded" {
		respondWithError(w, 204, fmt.Sprintf("%v", err))
		return
	}

	userID, err := uuid.Parse(reqData.Data.UserID)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("%v", err))
		return
	}

	if err := cfg.dbQueries.UpgradeUser(r.Context(), userID); err != nil {
		respondWithError(w, 404, fmt.Sprintf("%v", err))
		return
	}

	respondWithJson(w, 204, nil)

}
