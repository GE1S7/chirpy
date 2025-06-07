package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, req *http.Request) {

	header := w.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")

	cfg.fileserverHits.Swap(0)

	w.WriteHeader(http.StatusOK)
	msg := []byte(fmt.Sprintf("fileserverHits has been reset to 0"))
	w.Write(msg)
}
