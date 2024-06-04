package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Logger struct {
	infoLoggerFile  *log.Logger
	warnLoggerFile  *log.Logger
	errorLoggerFile *log.Logger
	infoLoggerTerm  *log.Logger
	warnLoggerTerm  *log.Logger
	errorLoggerTerm *log.Logger
	file            *os.File
}

var (
	instance *Logger
	once     sync.Once
)

func NewLogger(filePath string) *Logger {
	once.Do(func() {
		logDir := filepath.Dir(filePath)
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			log.Fatalf("Failed to create log directory: %s", err)
		}

		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file: %s", err)
		}

		instance = &Logger{
			infoLoggerFile:  log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
			warnLoggerFile:  log.New(file, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
			errorLoggerFile: log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
			infoLoggerTerm:  log.New(os.Stdout, "", 0),
			warnLoggerTerm:  log.New(os.Stdout, "", 0),
			errorLoggerTerm: log.New(os.Stdout, "", 0),
			file:            file,
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

func formatLogMessage(level, msg string) string {
	return fmt.Sprintf("[%s] [%s] %s", time.Now().Format("2006-01-02 15:04:05"), level, msg)
}

func formatLogMessageColor(level, color, msg string) string {
	return fmt.Sprintf("%s[%s] [%s] %s%s", color, time.Now().Format("2006-01-02 15:04:05"), level, msg, colorReset)
}

func (l *Logger) Info(msg string) {
	formattedMsg := formatLogMessage("INFO", msg)
	formattedMsgColor := formatLogMessageColor("INFO", colorGreen, msg)
	l.infoLoggerFile.Println(formattedMsg)
	l.infoLoggerTerm.Println(formattedMsgColor)
}

func (l *Logger) Warn(msg string) {
	formattedMsg := formatLogMessage("WARN", msg)
	formattedMsgColor := formatLogMessageColor("WARN", colorYellow, msg)
	l.warnLoggerFile.Println(formattedMsg)
	l.warnLoggerTerm.Println(formattedMsgColor)
}

func (l *Logger) Error(msg string) {
	formattedMsg := formatLogMessage("ERROR", msg)
	formattedMsgColor := formatLogMessageColor("ERROR", colorRed, msg)
	l.errorLoggerFile.Println(formattedMsg)
	l.errorLoggerTerm.Println(formattedMsgColor)
}

func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		log.Fatalf("Failed to close log file: %s", err)
	}
}
