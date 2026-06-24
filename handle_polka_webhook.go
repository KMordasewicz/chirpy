package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/KMordasewicz/chirpy/internal/auth"
)

func (cfg *apiConfig) polkaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetAPIKey(r.Header)
	if err != nil || token != cfg.polkaKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	msg, err := decodeMsg[msgPolkaWebhook](r)
	if err != nil {
		log.Printf("Couldn't decode Polka Webhook messsage: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if msg.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.dbQueires.UpgradeUserToRed(r.Context(), msg.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sendError(w, 404, "User not found")
			return
		} else {
			log.Printf("Couldn't retrive refresh token from db: %v\n", err)
			sendError(w, 500, "Couldn't verify refresh token")
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
