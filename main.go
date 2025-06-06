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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		fmt.Println("add 1")
		next.ServeHTTP(w, r)
	})

}

func healthHandler(w http.ResponseWriter, req *http.Request) {

	header := w.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)
	msg := []byte("OK")
	w.Write(msg)
}

func (cfg *apiConfig) reqNumHandler(w http.ResponseWriter, req *http.Request) {
	header := w.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)
	msg := []byte(fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load()))
	w.Write(msg)

}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, req *http.Request) {

	header := w.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")

	cfg.fileserverHits.Swap(0)

	w.WriteHeader(http.StatusOK)
	msg := []byte(fmt.Sprintf("fileserverHits has been reset to 0"))
	w.Write(msg)
}

func main() {
	cfg := apiConfig{}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fsh := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app", cfg.middlewareMetricsInc(fsh))

	mux.HandleFunc("GET /healthz", healthHandler)
	mux.HandleFunc("GET /metrics", cfg.reqNumHandler)
	mux.HandleFunc("POST /reset", cfg.resetHandler)

	log.Fatal(server.ListenAndServe())

}
