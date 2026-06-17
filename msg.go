package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type msgRecive interface {
	msgChirps | msgUsers
}

type responseSend interface {
	responseError | responseChirps | responseUsers | []responseChirps
}

type msgChirps struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type msgUsers struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type responseError struct {
	Error string `json:"error"`
}

type responseChirps struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type responseUsers struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func encodeMsg[R responseSend](m R, code int, w http.ResponseWriter) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func decodeMsg[M msgRecive](r *http.Request) (M, error) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var msg M
	if err := decoder.Decode(&msg); err != nil {
		return msg, err
	}
	return msg, nil
}

func sendError(w http.ResponseWriter, code int, msg string) {
	if msg == "" {
		msg = "Something went wrong"
	}
	resErr := responseError{Error: msg}
	err := encodeMsg(resErr, code, w)
	if err != nil {
		log.Printf("Couldn't encode error respone msg: %v\n", err)
	}
}
