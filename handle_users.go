package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) usersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "appliation/json")
	msg, err := decodeMsg[msgUsers](r)
	if err != nil {
		sendError(w, 400, "")
		return
	}
	users, err := cfg.dbQueires.CreateUser(r.Context(), msg.Email)
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
