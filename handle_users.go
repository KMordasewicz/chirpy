package main

import (
	"log"
	"net/http"

	"github.com/KMordasewicz/chirpy/internal/auth"
	"github.com/KMordasewicz/chirpy/internal/database"
)

func (cfg *apiConfig) usersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "appliation/json")
	msg, err := decodeMsg[msgUsers](r)
	if err != nil {
		sendError(w, 400, "")
		return
	}
	hashed_password, err := auth.HashPassword(msg.Password)
	if err != nil {
		log.Printf("Couldn't hash password: %v", err)
		sendError(w, 500, "Unable to process password")
		return
	}
	users, err := cfg.dbQueires.CreateUser(r.Context(), database.CreateUserParams{
		Email:          msg.Email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		log.Printf("Couldn't add user: %v\n", err)
		sendError(w, 500, "Couldn't add user.")
		return
	}
	taggedUsers := responseUsers{
		ID:        users.ID,
		CreatedAt: users.CreatedAt,
		UpdatedAt: users.UpdatedAt,
		Email:     users.Email,
	}
	err = encodeMsg(taggedUsers, 201, w)
	if err != nil {
		log.Printf("Couldn't encode respone msg: %v\n", err)
	}
}
