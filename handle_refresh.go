package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/KMordasewicz/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Couldn't get auth token: %v\n", err)
		sendError(w, 401, "Unauthorized")
		return
	}

	refreshTokenResult, err := cfg.dbQueires.GetUserFromRefreshToken(r.Context(), refreshToken)
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

	token, err := auth.MakeJWT(refreshTokenResult.ID, cfg.jwtSignKey)
	if err != nil {
		sendError(w, 500, "Couldn't generate auth token for the user")
		return
	}

	respMsg := responseRefresh{
		Token: token,
	}
	err = encodeMsg[responseRefresh](respMsg, http.StatusOK, w)
	if err != nil {
		log.Printf("Couldn't encode response msg: %v\n", err)
	}
}
