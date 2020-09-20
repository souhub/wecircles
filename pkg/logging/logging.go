package logging

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

var logFile *os.File

func init() {
	logFile, _ = os.OpenFile("./log.json", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(io.MultiWriter(logFile, os.Stdout))

	// Only log the info severity or above.
	log.SetLevel(log.WarnLevel)
}

func Info(msg string) {
	log.WithFields(log.Fields{
		"file": GetCurrentFile(),
		"line": GetCurrentFileLine(),
	}).Info(msg)
}

func Warn(msg string) {
	log.WithFields(log.Fields{
		"file": GetCurrentFile(),
		"line": GetCurrentFileLine(),
	}).Warn(msg)
}

func Fatal(msg string) {
	log.WithFields(log.Fields{
		"file": GetCurrentFile(),
		"line": GetCurrentFileLine(),
	}).Warn(msg)
}
