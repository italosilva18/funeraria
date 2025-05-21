package commands

import (
	"log"
	"os"

	"margem/robo/helpers"
)

var (
	logFile     *os.File
	logFilePath string
)

// SetLogFile sets the current log file handle and path used by commands.
func SetLogFile(f *os.File, path string) {
	logFile = f
	logFilePath = path
}

// ReopenLogFile closes the current log file, rotates if needed, and reopens it.
func ReopenLogFile() error {
	if logFile != nil {
		logFile.Close()
	}
	if err := helpers.RotateLogIfNeeded(logFilePath); err != nil {
		return err
	}
	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	logFile = f
	log.SetOutput(logFile)
	return nil
}
