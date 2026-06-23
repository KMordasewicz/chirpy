package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/KMordasewicz/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	platform       string
	dbQueires      *database.Queries
	jwtSignKey     string
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Couldn't load ENV vars: %v", err)
	}
	const serverPort = "8080"
	const rootDir = "."

	platform := os.Getenv("PLATFORM")
	jwtSignKey := os.Getenv("JWT_SIGN_KEY")
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Couldn't open database connetion: %v", err)
	}
	dbQueries := database.New(db)
	cfg := apiConfig{dbQueires: dbQueries, platform: platform, jwtSignKey: jwtSignKey}

	fileServerHandler := http.StripPrefix(
		"/app",
		http.FileServer(http.Dir(rootDir)),
	)

	serverMutex := http.NewServeMux()

	serverMutex.Handle("/app/", cfg.middlewareMetricsInc(fileServerHandler))

	serverMutex.HandleFunc("GET /api/healthz", healthzHandler)
	serverMutex.HandleFunc("POST /api/chirps", cfg.chirpsPostHandler)
	serverMutex.HandleFunc("GET /api/chirps", cfg.chirpsGetHandler)
	serverMutex.HandleFunc("GET /api/chirps/{chirpID}", cfg.chirpGetHandler)
	serverMutex.HandleFunc("POST /api/users", cfg.usersHandler)
	serverMutex.HandleFunc("PUT /api/users", cfg.updateUserHandler)
	serverMutex.HandleFunc("POST /api/login", cfg.loginHandler)
	serverMutex.HandleFunc("POST /api/refresh", cfg.refreshHandler)
	serverMutex.HandleFunc("POST /api/revoke", cfg.revokeHandler)

	serverMutex.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	serverMutex.HandleFunc("POST /admin/reset", cfg.resetHandler)

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
