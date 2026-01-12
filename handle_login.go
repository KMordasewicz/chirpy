package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/KMordasewicz/chirpy/internal/auth"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	msg, err := decodeMsg[msgUsers](r)
	if err != nil {
		sendError(w, 400, "")
		return
	}
	user, err := cfg.dbQueires.GetUserByEmail(r.Context(), msg.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user data for: %v\n", msg.Email)
			sendError(w, 401, "Incorrect email or password")
		} else {
			log.Printf("Coudn't get user data due, to: %v", err)
			sendError(w, 500, "")
		}
		return
	}
	match, err := auth.CheckPasswordHash(msg.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Couldn't match password: %v\n", err)
		sendError(w, 500, "Couldn't authenticate user")
		return
	}
	if !match {
		sendError(w, 401, "Incorrect email or password")
		return
	}
	res := responseUsers{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	if err = encodeMsg(res, 200, w); err != nil {
		log.Printf("Couldn't encode response: %v\n", err)
	}
}
