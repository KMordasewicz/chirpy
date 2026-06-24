package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const TokenTypeAccess TokenType = "chirpy-access"

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&jwt.RegisteredClaims{
			Issuer:    string(TokenTypeAccess),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			Subject:   userID.String(),
		},
	)

	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (any, error) {
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return uuid.Nil, err
	}
	if !token.Valid {
		return uuid.Nil, errors.New("token is invalid")
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, fmt.Errorf("incorrect issuer: %v", issuer)
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.Parse(subject)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	bearer := headers.Get("Authorization")
	if bearer == "" {
		return bearer, errors.New("no Authorization header")
	}
	token, is_bearer := strings.CutPrefix(bearer, "Bearer ")
	if !is_bearer {
		return "", errors.New("no Bearer in authorization header value")
	}
	return token, nil
}

func MakeRefreshToken() string {
	bytesBuffer := make([]byte, 32)
	rand.Read(bytesBuffer)

	return hex.EncodeToString(bytesBuffer)
}

func GetAPIKey(headers http.Header) (string, error) {
	apikey := headers.Get("Authorization")
	if apikey == "" {
		return apikey, errors.New("no Authorization header")
	}
	token, is_apikey := strings.CutPrefix(apikey, "ApiKey ")
	if !is_apikey {
		return "", errors.New("no ApiKey in authorization header value")
	}
	return token, nil
}
