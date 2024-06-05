package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"pc_booking_account_system/internal/data"
	"testing"
	"time"
)

type UsersTestSuite struct {
	suite.Suite
	app *application
	ts  *httptest.Server
}

func TestUsersSuite(t *testing.T) {
	suite.Run(t, new(UsersTestSuite))
}

func (s *UsersTestSuite) SetupSuite() {
	cfg := config{
		port: 4000,
		env:  "test",
		db: struct {
			dsn          string
			maxOpenConns int
			maxIdleConns int
			maxIdleTime  time.Duration
		}{
			dsn:          "postgres://default:fx1UvCaYwB8Q@ep-shy-cake-a2k2eyk6.eu-central-1.aws.neon.tech:5432/verceldb?sslmode=require",
			maxOpenConns: 25,
			maxIdleConns: 25,
			maxIdleTime:  15 * time.Minute, // Using time.Duration for maxIdleTime
		},
		limiter: struct {
			rps     float64
			burst   int
			enabled bool
		}{
			rps:     2,
			burst:   4,
			enabled: true,
		},
	}

	db, err := openDB(cfg)
	if err != nil {
		fmt.Println("err opening db")
	}
	//defer db.Close()

	models := data.NewModels(db)

	// test app
	s.app = &application{
		config: cfg,
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
		models: models,
	}

	// test server
	s.ts = httptest.NewServer(s.app.routes())
}

func (s *UsersTestSuite) TestRegisterUserHandler() {
	body := map[string]string{
		"fname":     "integration test",
		"sname":     "integration test",
		"email":     "integration test12@example.com",
		"password":  "inttest",
		"user-role": "test",
	}

	jsonBody, err := json.Marshal(body)
	s.Require().NoError(err, "failed to marshal request body")

	res, err := http.Post(s.ts.URL+"/pc_booking/user/register/", "application/json", bytes.NewBuffer(jsonBody))
	s.Require().NoError(err, "failed to make request")
	defer res.Body.Close()

	s.Require().Equal(http.StatusAccepted, res.StatusCode, "unexpected status code")
}

func (s *UsersTestSuite) TestGetUserByEmailHandler() {
	requestBody := `{"email": "test@example.com"}`

	req, err := http.NewRequest("GET", s.ts.URL+"/pc_booking/user/", bytes.NewBufferString(requestBody))
	s.Require().NoError(err, "failed to create request")

	res, err := http.DefaultClient.Do(req)
	s.Require().NoError(err, "failed to make request")
	defer res.Body.Close()

	s.Require().Equal(http.StatusOK, res.StatusCode, "unexpected status code")

	body, err := ioutil.ReadAll(res.Body)
	s.Require().NoError(err, "failed to read response body")

	fmt.Println(string(body))
}

func (s *UsersTestSuite) TestGetAllUsersHandler() {
	res, err := http.Get(s.ts.URL + "/pc_booking/users/all/")
	s.Require().NoError(err, "failed to make request")
	defer res.Body.Close()

	s.Require().Equal(http.StatusOK, res.StatusCode, "unexpected status code")

	body, err := ioutil.ReadAll(res.Body)
	s.Require().NoError(err, "failed to read response body")

	fmt.Println(string(body))
}

func (s *UsersTestSuite) TestDeleteByEmailUserHandler() {
	requestBody := `{"email": "test@example.com"}`

	req, err := http.NewRequest("DELETE", s.ts.URL+"/pc_booking/user/delete/", bytes.NewBufferString(requestBody))
	s.Require().NoError(err, "failed to create request")

	res, err := http.DefaultClient.Do(req)
	s.Require().NoError(err, "failed to make request")
	defer res.Body.Close()

	s.Require().Equal(http.StatusOK, res.StatusCode, "unexpected status code")

	body, err := ioutil.ReadAll(res.Body)
	s.Require().NoError(err, "failed to read response body")

	fmt.Println(string(body))
}

func (s *UsersTestSuite) TestActivateUserHandler() {
	requestBody := `{"token": "your_activation_token"}`

	req, err := http.NewRequest("PUT", s.ts.URL+"/pc_booking/users/activated", bytes.NewBufferString(requestBody))
	s.Require().NoError(err, "failed to create request")

	res, err := http.DefaultClient.Do(req)
	s.Require().NoError(err, "failed to make request")
	defer res.Body.Close()

	s.Require().Equal(http.StatusOK, res.StatusCode, "unexpected status code")

	body, err := ioutil.ReadAll(res.Body)
	s.Require().NoError(err, "failed to read response body")

	fmt.Println(string(body))
}
