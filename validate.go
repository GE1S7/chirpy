package main

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
)

func validateChirp(w http.ResponseWriter, req *http.Request, body string) (string, error) {

	if len(body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return "", fmt.Errorf("Chirp is too long")
	}

	return cleanOrgBody(body), nil

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
