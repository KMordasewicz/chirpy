package auth

import (
	"fmt"
	"net/http"
	"testing"
	"testing/synctest"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	testCases := []struct {
		password      string
		inputPassword string
		expectedMatch bool
	}{
		{
			password:      "pa$$word",
			inputPassword: "pa$$word",
			expectedMatch: true,
		},
		{
			password:      "admin123",
			inputPassword: "admin124",
			expectedMatch: false,
		},
	}
	for _, c := range testCases {
		hash, err := HashPassword(c.password)
		if err != nil {
			t.Fatalf("Error while hashing password: %v\n", err)
		}
		match, err := CheckPasswordHash(c.inputPassword, hash)
		if err != nil {
			t.Fatalf("Error while matching password: %v\n", err)
		}
		if c.expectedMatch != match {
			t.Errorf("Incorrect password validation for password: %v and input %v", c.password, c.inputPassword)
		}
	}
}

func TestMakeJWT(t *testing.T) {
	longDuration, _ := time.ParseDuration("30m")
	uuid1, _ := uuid.Parse("402d8009-dec5-4e32-b585-6db9e77776c3")
	uuid2, _ := uuid.Parse("8519a403-42fc-4923-8719-e8766e60c80c")
	testCases := []struct {
		name         string
		userID       uuid.UUID
		tokenSecret  string
		exipiersIn   time.Duration
		expectedFail bool
	}{
		{
			name:         "Good path",
			userID:       uuid1,
			tokenSecret:  "pa$$word",
			exipiersIn:   longDuration,
			expectedFail: false,
		},
		{
			name:         "Good path 2",
			userID:       uuid2,
			tokenSecret:  "😀", // Why does it work!?
			exipiersIn:   longDuration,
			expectedFail: false,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			_, err := MakeJWT(c.userID, c.tokenSecret, c.exipiersIn)
			if (err != nil) != c.expectedFail {
				fmt.Printf("MakeJWT() error = %v, expectedFail = %v\n", err, c.expectedFail)
				return
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	longDuration, _ := time.ParseDuration("30m")
	shortDuration, _ := time.ParseDuration("1s")
	uuid1, _ := uuid.Parse("402d8009-dec5-4e32-b585-6db9e77776c3")
	uuid2, _ := uuid.Parse("8519a403-42fc-4923-8719-e8766e60c80c")
	testCases := []struct {
		name              string
		userID            uuid.UUID
		tokenSecret       string
		exipiersIn        time.Duration
		passCorrectSecret bool
		expectedFail      bool
	}{
		{
			name:              "Good path",
			userID:            uuid1,
			tokenSecret:       "pa$$word",
			exipiersIn:        longDuration,
			passCorrectSecret: true,
			expectedFail:      false,
		},
		{
			name:              "Wrong password",
			userID:            uuid2,
			tokenSecret:       "test123",
			exipiersIn:        longDuration,
			passCorrectSecret: false,
			expectedFail:      true,
		},
		{
			name:              "Expired token",
			userID:            uuid1,
			tokenSecret:       "ha$hword",
			exipiersIn:        shortDuration,
			passCorrectSecret: true,
			expectedFail:      true,
		},
		{
			name:              "Wrong password and expired token",
			userID:            uuid2,
			tokenSecret:       "321test",
			exipiersIn:        shortDuration,
			passCorrectSecret: false,
			expectedFail:      true,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				token, err := MakeJWT(c.userID, c.tokenSecret, c.exipiersIn)
				if err != nil {
					t.Errorf("Token creation failed:\n\tuserID: %v\n\ttokenSecret: %v\n\texpiresIn: %v\nWith error:\n%v", c.userID, c.tokenSecret, c.exipiersIn, err)
					return
				}

				var secret string
				if c.passCorrectSecret {
					secret = c.tokenSecret
				} else {
					secret = "incorrectSecret123$"
				}

				time.Sleep(3 * time.Second)

				uuid, err := ValidateJWT(token, secret)
				if (err != nil) != c.expectedFail {
					fmt.Printf("ValidateJWT() error = %v, expectedFail = %v\n", err, c.expectedFail)
					return
				}

				if !c.expectedFail && uuid != c.userID {
					t.Errorf("uuid mismatch, got: %v instead: %v\n", uuid, c.userID)
				}
			})
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	testCases := []struct {
		name         string
		headers      http.Header
		token        string
		expectedFail bool
	}{
		{
			name:         "Good path",
			headers:      http.Header{"Authorization": []string{"Bearer 1234"}},
			token:        "1234",
			expectedFail: false,
		},
		{
			name:         "No key in headers",
			headers:      http.Header{},
			token:        "",
			expectedFail: true,
		},
		{
			name:         "No bearer",
			headers:      http.Header{"Authorization": []string{"Basic 1234"}},
			token:        "",
			expectedFail: true,
		},
	}
	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			token, err := GetBearerToken(c.headers)
			if (err != nil) != c.expectedFail {
				t.Errorf("GetBearerToken() error = %v, expectedFail = %v\n", err, c.expectedFail)
				return
			}
			if token != c.token {
				t.Errorf("token missmatch got: %v, expected: %v", token, c.token)
				return
			}
		})
	}
}
