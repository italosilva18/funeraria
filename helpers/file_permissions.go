package helpers

import (
	"os"
	"path/filepath"
)

// EnsureLogFile creates the log file and its parent directory if needed
// and sets its permissions to allow read/write for all users.
func EnsureLogFile(logPath string) error {
	dir := filepath.Dir(logPath)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}
	f, err := os.OpenFile(logPath, os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	f.Close()
	return os.Chmod(logPath, 0666)
}
