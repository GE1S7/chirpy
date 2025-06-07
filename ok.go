package main

import "net/http"

func healthHandler(w http.ResponseWriter, req *http.Request) {

	header := w.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)
	msg := []byte("OK")
	w.Write(msg)
}
