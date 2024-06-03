package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// Logger - структура нашего логгера
type Logger struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	file        *os.File
}

var (
	instance *Logger
	once     sync.Once
)

func NewLogger(filePath string) *Logger {
	once.Do(func() {
		// Ensure the directory exists
		logDir := filepath.Dir(filePath)
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			log.Fatalf("Failed to create log directory: %s", err)
		}

		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file: %s", err)
		}

		instance = &Logger{
			infoLogger:  log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
			warnLogger:  log.New(file, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
			errorLogger: log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
			file:        file,
		}
	})
	return instance
}

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

func (l *Logger) Info(msg string) {
	l.infoLogger.Println(msg)
	fmt.Printf("%sINFO: %s%s\n", colorGreen, msg, colorReset)
}

func (l *Logger) Warn(msg string) {
	l.warnLogger.Println(msg)
	fmt.Printf("%sWARN: %s%s\n", colorYellow, msg, colorReset)
}

func (l *Logger) Error(msg string) {
	l.errorLogger.Println(msg)
	fmt.Printf("%sERROR: %s%s\n", colorRed, msg, colorReset)
}

func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		log.Fatalf("Failed to close log file: %s", err)
	}
}
