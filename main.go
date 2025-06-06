package main

import (
	"fmt"
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

func healthHandler(w http.ResponseWriter, req *http.Request) {

	header := w.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")

	msg := []byte("OK")

	w.WriteHeader(http.StatusOK)
	w.Write(msg)
}

func (cfg *apiConfig) reqNumHandler(w http.ResponseWriter, req *http.Request) {

	header := w.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")

	msg := []byte(fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load()))

	w.WriteHeader(http.StatusOK)
	w.Write(msg)

}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, req *http.Request) {

	header := w.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")

	cfg.fileserverHits.Store(0)

	msg := []byte(fmt.Sprintf("fileserverHits has been reset to 0"))

	w.WriteHeader(http.StatusOK)
	w.Write(msg)
}

func main() {
	cfg := apiConfig{}
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	mux.HandleFunc("/healthz", healthHandler)
	mux.HandleFunc("/metrics", cfg.reqNumHandler)
	mux.HandleFunc("/reset", cfg.resetHandler)

	fsh := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))

	mux.Handle("/app/", cfg.middlewareMetricsInc(fsh))

	log.Fatal(server.ListenAndServe())

}
