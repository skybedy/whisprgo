package logx

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	localLogName    = "whisprgo.log"
	fallbackLogPath = ".local/state/whisprgo/whisprgo.log"
)

type Logger struct {
	path string
}

func New() (*Logger, error) {
	// Primary location: log file directly in current app directory.
	if wd, err := os.Getwd(); err == nil {
		localPath := filepath.Join(wd, localLogName)
		if f, err := os.OpenFile(localPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644); err == nil {
			_ = f.Close()
			return &Logger{path: localPath}, nil
		}
	}

	// Fallback: user state directory if local path is unavailable.
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to resolve home directory for logging: %w", err)
	}

	path := filepath.Join(home, fallbackLogPath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	return &Logger{path: path}, nil
}

func (l *Logger) Path() string {
	return l.path
}

func (l *Logger) Info(msg string) {
	l.write("INFO", msg)
}

func (l *Logger) Error(msg string) {
	l.write("ERROR", msg)
}

func (l *Logger) write(level string, msg string) {
	f, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()

	ts := time.Now().Format(time.RFC3339)
	_, _ = fmt.Fprintf(f, "%s [%s] %s\n", ts, level, msg)
}
