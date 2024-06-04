package main

import (
	"booking/internal/config"
	db2 "booking/internal/db"
	"booking/internal/jwt"
	"booking/internal/logger"
	"booking/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

type application struct {
	logger  *logger.Logger
	router  *gin.Engine
	db      *gorm.DB
	config  *config.Config
	jwtUtil *jwt.JWTUtil
}

func (app *application) initializeRoutes() {
	app.router.POST("/book", app.authMiddleware, app.bookComputer)
	app.router.GET("/available", app.getAvailableComputers)
	app.router.DELETE("/cancel/:id", app.authMiddleware, app.cancelBooking)
}

func newApplication() (*application, error) {

	cfg := config.LoadConfig()

	myLogger := logger.NewLogger(cfg.LogPath)

	myLogger.Info("Config loaded from env")

	dbconn, err := db2.InitDB(*cfg)
	if err != nil {
		myLogger.Error("Error initializing DB")
		return nil, err
	}

	if err := dbconn.AutoMigrate(&models.Computer{}, &models.Booking{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	db2.SeedComputers(dbconn)

	myLogger.Info("DB initialized")

	gin.SetMode(cfg.GinMode)

	router := gin.New()

	app := &application{
		logger:  myLogger,
		router:  router,
		db:      dbconn,
		jwtUtil: jwt.NewJWTUtil(cfg.Secret),
		config:  cfg,
	}

	app.initializeRoutes()

	myLogger.Info("routes initialized")
	return app, nil
}

func main() {
	app, err := newApplication()
	if err != nil {
		log.Fatalf("Could not initialize application: %s", err)
	}

	defer app.logger.Close()

	app.logger.Info("Starting booking service")
	if err := app.router.Run(":8080"); err != nil {
		app.logger.Error(fmt.Sprintf("Could not start server: %s", err))
	}
}
