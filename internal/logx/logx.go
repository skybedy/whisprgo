package logx

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	logDirName  = "whisprgo"
	logFileName = "whisprgo.log"
)

type Logger struct {
	path string
}

func New() (*Logger, error) {
	return NewWithPath("")
}

func NewWithPath(customPath string) (*Logger, error) {
	if customPath != "" {
		path := filepath.Clean(customPath)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}
		return &Logger{path: path}, nil
	}

	stateDir := os.Getenv("XDG_STATE_HOME")
	if stateDir == "" {
		home, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return nil, fmt.Errorf("failed to resolve log directory: %w", homeErr)
		}
		stateDir = filepath.Join(home, ".local", "state")
	}

	path := filepath.Join(stateDir, logDirName, logFileName)
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
