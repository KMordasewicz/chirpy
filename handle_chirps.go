package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/KMordasewicz/chirpy/internal/database"
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

func (cfg *apiConfig) chirpsHandler(w http.ResponseWriter, r *http.Request) {
	msg, err := decodeMsg[msgChirps](r)
	if err != nil {
		resErr := responseError{Error: "Something went wrong."}
		err := encodeMsg(resErr, 400, w)
		if err != nil {
			log.Printf("Couldn't encode respone msg: %v\n", err)
		}
		return
	}

	if len(msg.Body) > 140 {
		resErr := responseError{Error: "Chirp is too long"}
		err := encodeMsg(resErr, 400, w)
		if err != nil {
			log.Printf("Couldn't encode Chirpy too long msg: %v\n", err)
		}
		return
	}

	chirp, err := cfg.dbQueires.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanMsg(msg.Body),
		UserID: msg.UserID,
	})
	if err != nil {
		log.Printf("Couldn't add chrip: %v\n", err)
		resErr := responseError{Error: "Something went wrong."}
		err := encodeMsg(resErr, 400, w)
		if err != nil {
			log.Printf("Couldn't encode respone msg: %v\n", err)
		}
		return
	}

	res := responseChirps{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	err = encodeMsg(res, 201, w)
	if err != nil {
		log.Printf("Couldn't encode OK msg: %v\n", err)
	}
}
