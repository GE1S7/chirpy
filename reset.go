package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform == "dev" {

		header := w.Header()
		header.Set("Content-Type", "text/plain; charset=utf-8")

		cfg.fileserverHits.Swap(0)

		err := cfg.dbQueries.DeleteUsers(r.Context())
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
		msg := []byte(fmt.Sprintf("fileserverHits has been reset to 0\nAll users has been deleted"))
		w.Write(msg)

	} else {
		w.WriteHeader(403)
	}

}
