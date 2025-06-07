package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func validateHandler(w http.ResponseWriter, req *http.Request) {
	type resp struct {
		Resp string `json:"cleaned_body"`
	}

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {

		respondWithError(w, 500, "Error decoding request")

	} else if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
	} else {

		respValid := resp{
			Resp: cleanOrgBody(params.Body),
		}

		respondWithJson(w, 200, respValid)

	}

}

func cleanOrgBody(body string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(body, " ")

	for i, w := range words {
		lowerWord := strings.ToLower(w)
		if slices.Contains(badWords, lowerWord) {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
