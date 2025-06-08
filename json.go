package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type fail struct {
		Error string `json:"error"`
	}

	respFail := fail{
		Error: msg,
	}

	respondWithJson(w, code, respFail)

}
