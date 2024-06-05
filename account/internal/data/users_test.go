package data

import (
	"database/sql"
	"fmt"
	"testing"
)

var userModel UserModel

func TestGetForToken(t *testing.T) {
	var tokenPlaintext string
	tokenPlaintext = ""
	setUserTestsVars()
	//"authentication"
	user, err := userModel.GetForToken("authentication", tokenPlaintext)
	if err != nil {
		t.Errorf("Error getting user by token")
	}
	if user == nil {
		t.Errorf("User is null, but not supposed to")
	} else {
		t.Logf("User id: %d", user.ID)
		t.Logf("User name: %s %s", user.SName, user.FName)
	}
}

func TestGetByEmail(t *testing.T) {
	setUserTestsVars()
	var email string
	email = "xakercool33@gmail.com"

	user, err := userModel.GetByEmail(email)
	if err != nil {
		t.Errorf("Error getting user by email: %e", err)
	}
	if user == nil {
		t.Errorf("User is null, but not supposed to")
	} else {
		t.Logf("User id: %d", user.ID)
		t.Logf("User name: %s %s", user.SName, user.FName)
	}
}

func TestPwdMatches(t *testing.T) {
	setUserTestsVars()
	var plaintextPassword string
	plaintextPassword = "admin"

	user, err := userModel.Get("afc2a7fc-f62f-4d33-82d7-959b34b54c43")

	matches, err := user.Password.Matches(plaintextPassword)
	if err != nil {
		t.Errorf("Error while matching passwords: %e", err)
	}
	if !matches {
		t.Errorf("Passwords don't match!")
	} else {
		t.Log("Passwords matched!")
	}
}

func setUserTestsVars() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/kh.takenovDB?sslmode=disable")
	if err != nil {
		fmt.Println("Error open db")
	}

	userModel = UserModel{DB: db}
}
