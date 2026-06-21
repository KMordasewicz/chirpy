package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/KMordasewicz/chirpy/internal/auth"
)

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	// if r.Body != nil {
	// 	sendError(w, 400, "no body accepted")
	// 	return
	// }

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Couldn't get auth token: %v\n", err)
		sendError(w, 401, "Unauthorized")
		return
	}

	_, err = cfg.dbQueires.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sendError(w, 401, "Invalid refresh token")
			return
		} else {
			log.Printf("Couldn't retrive refresh token from db: %v\n", err)
			sendError(w, 500, "Couldn't verify refresh token")
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
