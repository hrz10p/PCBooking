package main

import (
	db2 "booking/internal/db"
	"booking/internal/jwt"
	"booking/internal/logger"
	"booking/internal/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func setupTestApp() (*application, error) {
	gin.SetMode(gin.TestMode)
	logger := logger.NewLogger("../logs/test.log")

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger2.Default.LogMode(logger2.Silent),
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Computer{}, &models.Booking{})
	if err != nil {
		return nil, err
	}

	jwtUtil := jwt.NewJWTUtil("test-secret")

	app := &application{
		router:  gin.New(),
		db:      db,
		logger:  logger,
		jwtUtil: jwtUtil,
	}

	app.initializeRoutes()

	return app, nil
}

func TestBookComputer(t *testing.T) {
	app, err := setupTestApp()
	assert.NoError(t, err)

	computer := models.Computer{Number: 1, Status: "available"}
	app.db.Create(&computer)

	token, err := app.jwtUtil.GenerateToken("user123", "user", "email@mail.com")
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/book", strings.NewReader(`{"computer_id": 1, "start_time": "2024-06-01T10:00:00Z", "end_time": "2024-06-01T12:00:00Z"}`))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Computer booked successfully")
}

func TestGetAvailableComputers(t *testing.T) {
	app, err := setupTestApp()
	assert.NoError(t, err)

	db2.SeedComputers(app.db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/available", nil)

	var computers []models.Computer

	err = json.Unmarshal(w.Body.Bytes(), &computers)
	if err != nil {
		return
	}
	app.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, len(computers), 6)
}

func TestCancelBooking(t *testing.T) {
	app, err := setupTestApp()
	assert.NoError(t, err)

	computer := models.Computer{Number: 1, Status: "available"}
	app.db.Create(&computer)

	stTime, err := time.Parse(time.DateTime, "2024-06-01 10:00:00")
	assert.NoError(t, err)

	endTime, err := time.Parse(time.DateTime, "2024-06-01 12:00:00")
	assert.NoError(t, err)

	booking := models.Booking{
		UserID:     "user123",
		ComputerID: computer.ID,
		StartTime:  stTime,
		EndTime:    endTime,
	}
	app.db.Create(&booking)

	token, err := app.jwtUtil.GenerateToken("user123", "user", "email@mail.com")
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/cancel/"+strconv.Itoa(int(booking.ID)), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Booking cancelled")
}

func TestAdminCreateComputer(t *testing.T) {
	app, err := setupTestApp()
	assert.NoError(t, err)

	token, err := app.jwtUtil.GenerateToken("user123", "admin", "email@mail.com")
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/admin/computers", strings.NewReader(`{"number": 2, "status": "available"}`))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "available")
}

func TestAdminMiddleware(t *testing.T) {
	app, err := setupTestApp()
	assert.NoError(t, err)

	token, err := app.jwtUtil.GenerateToken("user123", "user", "email@mail.com")
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/admin/computers", strings.NewReader(`{"number": 2, "status": "available"}`))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Forbidden")
}
