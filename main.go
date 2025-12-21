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

func main() {
	const serverPort = "8080"
	const rootDir = "."
	cfg := apiConfig{}

	fileServerHandler := http.StripPrefix(
		"/app",
		http.FileServer(http.Dir(rootDir)),
	)

	serverMutex := http.NewServeMux()
	serverMutex.Handle("/app/", cfg.middlewareMetricsInc(fileServerHandler))
	serverMutex.HandleFunc("GET /api/healthz", healthzHandler)
	serverMutex.HandleFunc("GET /api/metrics", cfg.metricsHandler)
	serverMutex.HandleFunc("POST /api/reset", cfg.resetHandler)

	server := &http.Server{
		Addr:    ":" + serverPort,
		Handler: serverMutex,
	}

	fmt.Println("Starting Chirpy sever!")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server stopped, due to: %v\n", err)
	}
	fmt.Println("Server shutdown")
}
