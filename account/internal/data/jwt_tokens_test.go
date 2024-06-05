package data

import (
	"github.com/golang-jwt/jwt/v4"
	"testing"
)

var (
	claims        Claims
	jwtTokenModel JWTUtil
)

func TestGenerateToken(t *testing.T) {
	token, err := jwtTokenModel.GenerateToken("afc2a7fc-f62f-4d33-82d7-959b34b54c43", "admin", "xakercool33@gmail.com")
	if err != nil {
		t.Errorf("Error occured while generating jwt token: %e", err)
	}
	t.Logf("Token: %s", token)
}

func TestValidateToken(t *testing.T) {
	token, err := jwt.ParseWithClaims("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJhZmMyYTdmYy1mNjJmLTRkMzMtODJkNy05NTliMzRiNTRjNDMiLCJyb2xlIjoiYWRtaW4iLCJlbWFpbCI6Inhha2VyY29vbDMzQGdtYWlsLmNvbSIsImV4cCI6MTcxNzc4MDMzNH0.J20W7z5n2Yo63YJKLsYLN5JPcBrAicyZy2eVm4SK1oM", &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtTokenModel.secret), nil
	})

	if err != nil {
		t.Errorf("Error occured: %e", err)
	}

	_, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		t.Errorf("Err: invalid token: %e", err)
	}
}
