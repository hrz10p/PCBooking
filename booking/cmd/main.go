package main

import (
	"booking/internal/config"
	db2 "booking/internal/db"
	"booking/internal/logger"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// application - структура приложения, содержащая логгер, маршруты и соединение с базой данных
type application struct {
	logger *logger.Logger
	router *gin.Engine
	db     *gorm.DB
	config *config.Config
}

// initializeRoutes - функция для инициализации маршрутов
func (app *application) initializeRoutes() {
	app.router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Добавьте сюда другие маршруты
}

// newApplication - функция для создания нового экземпляра приложения
func newApplication() (*application, error) {

	cfg := config.LoadConfig()

	myLogger := logger.NewLogger(cfg.LogPath)

	myLogger.Info("Config loaded from env")

	db2.InitDB(*cfg)

	myLogger.Info("DB initialized")

	gin.SetMode(cfg.GinMode)

	router := gin.New()

	app := &application{
		logger: myLogger,
		router: router,
		db:     db2.DB,
		config: cfg,
	}

	app.initializeRoutes()

	myLogger.Info("routes initialized")
	return app, nil
}

func main() {
	// Создание нового приложения
	app, err := newApplication()
	if err != nil {
		log.Fatalf("Could not initialize application: %s", err)
	}

	defer app.logger.Close()

	// Запуск сервера
	app.logger.Info("Starting booking service")
	if err := app.router.Run(":8080"); err != nil {
		app.logger.Error(fmt.Sprintf("Could not start server: %s", err))
	}
}
