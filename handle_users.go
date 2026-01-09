package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) usersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "appliation/json")
	msg, err := decodeMsg[msgUsers](r)
	if err != nil {
		resErr := responseError{Error: "Something went wrong."}
		err := encodeMsg(resErr, 400, w)
		if err != nil {
			log.Printf("Couldn't dencode respone msg: %v\n", err)
		}
		return
	}
	users, err := cfg.dbQueires.CreateUser(r.Context(), msg.Email)
	if err != nil {
		resErr := responseError{Error: "Couldn't add user."}
		log.Printf("Couldn't add user: %v\n", err)
		err := encodeMsg(resErr, 500, w)
		if err != nil {
			log.Printf("Couldn't encode msg: %v\n", err)
		}
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
