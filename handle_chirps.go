package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/KMordasewicz/chirpy/internal/auth"
	"github.com/KMordasewicz/chirpy/internal/database"
	"github.com/google/uuid"
)

var profaineWords = [...]string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

func checkProfanity(word string) string {
	for _, v := range profaineWords {
		if strings.ToLower(word) == v {
			return "****"
		}
	}
	return word
}

func cleanMsg(msg string) string {
	words := strings.Split(msg, " ")
	cleanMsg := make([]string, 0, len(words))
	for _, word := range words {
		cleanMsg = append(cleanMsg, checkProfanity(word))
	}
	return strings.Join(cleanMsg, " ")
}

func (cfg *apiConfig) chirpsPostHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Couldn't get auth token: %v\n", err)
		sendError(w, 401, "Unauthorized")
		return
	}

	msg, err := decodeMsg[msgChirps](r)
	if err != nil {
		log.Printf("Couldn't decode message: %v\n", err)
		sendError(w, 400, "Incorrect message format")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSignKey)
	if err != nil {
		log.Printf("Invalid token user token: %v\n", err)
		sendError(w, 401, "Unauthorized")
		return
	}
	// if user != msg.UserID {
	// 	log.Printf("User missmatch msg user: %v, token: %v", msg.UserID, user)
	// 	sendError(w, 401, "Incorrect user authorization")
	// 	return
	// }

	if len(msg.Body) > 140 {
		sendError(w, 400, "Chirp is too long")
		return
	}

	chirp, err := cfg.dbQueires.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanMsg(msg.Body),
		UserID: userID,
	})
	if err != nil {
		log.Printf("Couldn't add chrip: %v\n", err)
		sendError(w, 400, "")
		return
	}

	res := responseChirps{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.ID,
	}
	err = encodeMsg(res, 201, w)
	if err != nil {
		log.Printf("Couldn't encode response msg: %v\n", err)
	}
}

func (cfg *apiConfig) chirpsGetHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.dbQueires.GetChirps(r.Context())
	if err != nil {
		log.Printf("Couldn't get chirps from db: %v\n", err)
		sendError(w, 500, "Unable to fetch chirps")
		return
	}
	chirpsTagged := make([]responseChirps, len(chirps))
	for i, chirp := range chirps {
		chirpsTagged[i] = responseChirps{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}
	err = encodeMsg(chirpsTagged, 200, w)
	if err != nil {
		log.Printf("Couldn't encode response msg: %v\n", err)
	}
}

func (cfg *apiConfig) chirpGetHandler(w http.ResponseWriter, r *http.Request) {
	chirpIDSting := r.PathValue("chirpID")
	if chirpIDSting == "" {
		log.Print("Failed to get path value for chirp id.\n")
		return
	}
	chirpID, err := uuid.Parse(chirpIDSting)
	if err != nil {
		log.Printf("Couldn't parse chirp id to uuid: %v\n", err)
		return
	}
	chirp, err := cfg.dbQueires.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No rows for id: %v", chirpID)
			sendError(w, 404, "Chirp not found")
		} else {
			log.Printf("Couldn't get chirp for id: %v from db: %v\n", chirpID, err)
			sendError(w, 500, "Unable to fetch chirp")
		}
		return
	}
	chirpTagged := responseChirps{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	err = encodeMsg(chirpTagged, 200, w)
	if err != nil {
		log.Printf("Couldn't encode response msg: %v\n", err)
	}
}
