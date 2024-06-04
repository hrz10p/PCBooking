package logger

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	filePath := "../../logs/test.log"
	_ = os.Remove(filePath) // ensure we start fresh

	log := NewLogger(filePath)
	defer log.Close()

	log.Info("Test log message")
	log.Warn("Test warn message")
	log.Error("Test error message")

	content, err := ioutil.ReadFile(filePath)
	assert.NoError(t, err)

	assert.Contains(t, string(content), "Test log message")
	assert.Contains(t, string(content), "Test warn message")
	assert.Contains(t, string(content), "Test error message")
}
