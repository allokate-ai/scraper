package logger

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/joho/godotenv"
)

type Logger struct {
	client *fluent.Fluent
}

func (l *Logger) Error(message string) {
	var data = map[string]string{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"message":   message,
		"level":     "ERROR",
	}
	l.client.Post("auth.ERROR", data)
}

func (l *Logger) Info(message string) {
	var data = map[string]string{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"message":   message,
		"level":     "INFO",
	}
	l.client.Post("auth.INFO", data)

}

func (l *Logger) Debug(message string) {
	var data = map[string]string{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"message":   message,
		"level":     "DEBUG",
	}
	l.client.Post("auth.DEBUG", data)

}

func (l *Logger) Trace(message string) {
	var data = map[string]string{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"message":   message,
		"level":     "TRACE",
	}
	l.client.Post("auth.TRACE", data)
}

func (l *Logger) Close() {
	l.client.Close()
}

var logger Logger

func Get() *Logger {
	godotenv.Load()
	port, _ := strconv.Atoi(os.Getenv("FLUENTD_PORT"))
	if logger.client == nil {
		client, err := fluent.New(fluent.Config{
			FluentHost: os.Getenv("FLUENTD_HOST"),
			FluentPort: port,
		})
		if err != nil {
			log.Fatal(err)
		}
		logger.client = client
	}
	return &logger
}
