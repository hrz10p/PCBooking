package data

import (
	"testing"
	"time"
)

var (
	model TokenModel
)

func TestGeneratingToken(t *testing.T) {
	userId := "afc2a7fc-f62f-4d33-82d7-959b34b54c43"
	ttl := 20 * time.Minute
	scope := "activation"

	token, err := generateToken(userId, ttl, scope)
	if err != nil {
		t.Errorf("Error while generating token: %e", err)
	}
	if len(token.Plaintext) > 26 {
		t.Errorf("Token len must not be more than 26, bt got 26")
	} else if len(token.Plaintext) == 0 {
		t.Errorf("Token len must be 26, bt got 0")
	} else {
		t.Logf("Token is good: %s", token.Plaintext)
	}
}
