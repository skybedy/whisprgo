package parakeet

import (
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"whisprgo/internal/config"
)

func TestValidateManagedConfigSuccess(t *testing.T) {
	cfg := validManagedConfig(t)
	if err := ValidateManagedConfig(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateManagedConfigMissingModelFile(t *testing.T) {
	cfg := validManagedConfig(t)
	if err := os.Remove(filepath.Join(cfg.ModelDir, "joiner.int8.onnx")); err != nil {
		t.Fatalf("remove: %v", err)
	}
	err := ValidateManagedConfig(cfg)
	if err == nil || !strings.Contains(err.Error(), "joiner.int8.onnx") {
		t.Fatalf("expected missing model file error, got %v", err)
	}
}

func TestCheckExternalEndpoint(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()

	addr := ln.Addr().String()
	if err := CheckExternalEndpoint("ws://"+addr, time.Second); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func validManagedConfig(t *testing.T) config.ParakeetConfig {
	t.Helper()
	dir := t.TempDir()
	binaryPath := filepath.Join(dir, "sherpa")
	if err := os.WriteFile(binaryPath, []byte("stub"), 0o755); err != nil {
		t.Fatalf("write binary: %v", err)
	}
	modelDir := filepath.Join(dir, "model")
	if err := os.MkdirAll(modelDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	for _, name := range []string{"tokens.txt", "encoder.int8.onnx", "decoder.int8.onnx", "joiner.int8.onnx"} {
		if err := os.WriteFile(filepath.Join(modelDir, name), []byte("stub"), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	return config.ParakeetConfig{
		Mode:                  "managed",
		Binary:                binaryPath,
		ModelDir:              modelDir,
		Host:                  "127.0.0.1",
		Port:                  6010,
		NumThreads:            4,
		StartupTimeoutSeconds: 10,
		RequestTimeoutSeconds: 120,
	}
}
