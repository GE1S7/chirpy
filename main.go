package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	_ = cfg.fileserverHits.Add(1)
	return next
}

func handler(w http.ResponseWriter, req *http.Request) {

	header := w.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)
	msg := []byte("OK")
	w.Write(msg)
}

func main() {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	mux.HandleFunc("/healthz", handler)

	fsHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))

	mux.Handle("/app/", fsHandler)

	log.Fatal(server.ListenAndServe())

}
