package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type msgRecive interface {
	msgChirps | msgUsers
}

type responseSend interface {
	responseError | responseChirps | responseUsers
}

type msgChirps struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type msgUsers struct {
	Email string `json:"email"`
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
