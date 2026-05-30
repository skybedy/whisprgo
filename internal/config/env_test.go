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

func TestUpsertAndRemoveKeyInEnvFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")

	if err := UpsertKeyInEnvFile(path, "OPENAI_API_KEY", "first"); err != nil {
		t.Fatalf("UpsertKeyInEnvFile create failed: %v", err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read .env failed: %v", err)
	}
	if string(raw) != "OPENAI_API_KEY=first\n" {
		t.Fatalf("unexpected content after create: %q", string(raw))
	}

	if err := UpsertKeyInEnvFile(path, "OPENAI_API_KEY", "second"); err != nil {
		t.Fatalf("UpsertKeyInEnvFile update failed: %v", err)
	}
	raw, err = os.ReadFile(path)
	if err != nil {
		t.Fatalf("read .env failed: %v", err)
	}
	if string(raw) != "OPENAI_API_KEY=second\n" {
		t.Fatalf("unexpected content after update: %q", string(raw))
	}

	if err := UpsertKeyInEnvFile(path, "OTHER", "x"); err != nil {
		t.Fatalf("UpsertKeyInEnvFile append failed: %v", err)
	}
	if err := RemoveKeyFromEnvFile(path, "OPENAI_API_KEY"); err != nil {
		t.Fatalf("RemoveKeyFromEnvFile failed: %v", err)
	}
	raw, err = os.ReadFile(path)
	if err != nil {
		t.Fatalf("read .env failed: %v", err)
	}
	if string(raw) != "OTHER=x\n" {
		t.Fatalf("unexpected content after remove: %q", string(raw))
	}
}
