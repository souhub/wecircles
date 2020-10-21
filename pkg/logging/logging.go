package logging

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

var logFile *os.File

func init() {
	logFile, _ = os.OpenFile("./log.json", os.O_WRONLY|os.O_CREATE, 0755)

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(io.MultiWriter(logFile, os.Stdout))

	// Only log the info severity or above.
	log.SetLevel(log.InfoLevel)
}

func Info(err error, fileName string, line int) {
	log.WithFields(log.Fields{
		"file": fileName,
		"line": line,
	}).Info(err)
}

func Warn(err error, fileName string, line int) {
	log.WithFields(log.Fields{
		"file": fileName,
		"line": line,
	}).Warn(err)
}

func Fatal(err error, fileName string, line int) {
	log.WithFields(log.Fields{
		"file": fileName,
		"line": line,
	}).Fatal(err)
}
