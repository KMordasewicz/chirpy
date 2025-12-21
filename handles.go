package main

import (
	"fmt"
	"log"
	"net/http"
)

func healthzHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	_, err := rw.Write([]byte("OK"))
	if err != nil {
		log.Fatalf("Could write message to healthz handler: %v\n", err)
	}
}

func (cfg *apiConfig) metricsHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(rw, "Hits: %v", cfg.fileserverHits.Load())
	if err != nil {
		log.Fatalf("Could write message to healthz handler: %v\n", err)
	}
}

func (cfg *apiConfig) resetHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	_, err := rw.Write([]byte("Hits reset"))
	if err != nil {
		log.Fatalf("Could write message to healthz handler: %v\n", err)
	}
}
