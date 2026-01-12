package auth

import (
	"testing"
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
