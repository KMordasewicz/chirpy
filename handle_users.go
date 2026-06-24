package main

import (
	"log"
	"net/http"

	"github.com/KMordasewicz/chirpy/internal/auth"
	"github.com/KMordasewicz/chirpy/internal/database"
)

func (cfg *apiConfig) usersHandler(w http.ResponseWriter, r *http.Request) {
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
		ID:          users.ID,
		CreatedAt:   users.CreatedAt,
		UpdatedAt:   users.UpdatedAt,
		Email:       users.Email,
		IsChirpyRed: users.IsChirpyRed,
	}
	err = encodeMsg(taggedUsers, 201, w)
	if err != nil {
		log.Printf("Couldn't encode respone msg: %v\n", err)
	}
}

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSignKey)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	msg, err := decodeMsg[msgUsers](r)
	if err != nil {
		sendError(w, http.StatusBadRequest, "")
		return
	}

	hashedPassword, err := auth.HashPassword(msg.Password)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Couldn't update password")
		return
	}

	userResult, err := cfg.dbQueires.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
		ID:             userID,
		Email:          msg.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Couldn't update password")
		return
	}

	err = encodeMsg(responseUsers{
		ID:          userResult.ID,
		CreatedAt:   userResult.CreatedAt,
		UpdatedAt:   userResult.UpdatedAt,
		Email:       userResult.Email,
		IsChirpyRed: userResult.IsChirpyRed,
	}, http.StatusOK, w)
	if err != nil {
		log.Printf("Couldn't respond to POST /api/users request: %v\n", err)
	}
}
