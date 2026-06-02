package app

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteForegroundSkipsParakeetInfo(t *testing.T) {
	var buf bytes.Buffer
	a := &App{foregroundLogs: true, foregroundWriter: &buf}

	a.writeForeground("INFO", nil, "parakeet stdout: noisy line")
	if buf.Len() != 0 {
		t.Fatalf("expected parakeet info log to be skipped, got %q", buf.String())
	}

	a.writeForeground("INFO", nil, "cleanup: enabled")
	if !strings.Contains(buf.String(), "cleanup: enabled") {
		t.Fatalf("expected non-parakeet info log to be written, got %q", buf.String())
	}
}

func TestWriteForegroundKeepsParakeetErrors(t *testing.T) {
	var buf bytes.Buffer
	a := &App{foregroundLogs: true, foregroundWriter: &buf}

	a.writeForeground("ERROR", nil, "parakeet start failed err=boom")
	if !strings.Contains(buf.String(), "parakeet start failed err=boom") {
		t.Fatalf("expected parakeet error log to be written, got %q", buf.String())
	}
}
