package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveSecretPrefersEnvironment(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("OPENAI_API_KEY", "from-env")

	cfgDir := filepath.Join(home, ".config", "whisprgo")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(cfgDir, ".env"), []byte("OPENAI_API_KEY=from-dotenv\n"), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	got := ResolveSecret("OPENAI_API_KEY")
	if got != "from-env" {
		t.Fatalf("expected env priority, got %q", got)
	}
}

func TestResolveSecretFallsBackToDotEnv(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("OPENAI_API_KEY", "")

	cfgDir := filepath.Join(home, ".config", "whisprgo")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(cfgDir, ".env"), []byte("export OPENAI_API_KEY='from-dotenv'\n"), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	got := ResolveSecret("OPENAI_API_KEY")
	if got != "from-dotenv" {
		t.Fatalf("expected dotenv fallback, got %q", got)
	}
}

func TestResolveSecretPrefersLocalDotEnvOverHomeDotEnv(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("OPENAI_API_KEY", "")

	cfgDir := filepath.Join(home, ".config", "whisprgo")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(cfgDir, ".env"), []byte("OPENAI_API_KEY=from-home-dotenv\n"), 0o644); err != nil {
		t.Fatalf("write home .env: %v", err)
	}

	projectDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(projectDir, ".env"), []byte("OPENAI_API_KEY=from-local-dotenv\n"), 0o644); err != nil {
		t.Fatalf("write local .env: %v", err)
	}

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() { _ = os.Chdir(origWD) }()
	if err := os.Chdir(projectDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	got := ResolveSecret("OPENAI_API_KEY")
	if got != "from-local-dotenv" {
		t.Fatalf("expected local .env priority, got %q", got)
	}
}
