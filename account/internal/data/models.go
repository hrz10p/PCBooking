package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrTokenExpired   = errors.New("token expired")
)

type Models struct {
	Users            UserModel
	ActivationTokens TokenModel
	JWTTokens        JWTUtil
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:            UserModel{DB: db},
		ActivationTokens: TokenModel{DB: db},
		JWTTokens:        JWTUtil{secret: "secret"},
	}
}
