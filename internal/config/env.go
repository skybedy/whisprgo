package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SecretSource string

const (
	SecretSourceEnvironment SecretSource = "environment variable"
	SecretSourceLocalEnv    SecretSource = "./.env"
	SecretSourceHomeEnv     SecretSource = "~/.config/whispergo/.env"
	SecretSourceMissing     SecretSource = "missing"
)

func EnvPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "~/.config/whispergo/.env"
	}
	return filepath.Join(home, ".config", "whispergo", ".env")
}

func ResolveSecret(key string) string {
	v, _ := ResolveSecretWithSource(key)
	return v
}

func ResolveSecretWithSource(key string) (string, SecretSource) {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v, SecretSourceEnvironment
	}

	paths := []struct {
		path   string
		source SecretSource
	}{
		{path: LocalEnvPath(), source: SecretSourceLocalEnv},
		{path: EnvPath(), source: SecretSourceHomeEnv},
	}
	for _, p := range paths {
		pairs, err := readDotEnv(p.path)
		if err != nil {
			continue
		}
		if v := strings.TrimSpace(pairs[key]); v != "" {
			return v, p.source
		}
	}

	return "", SecretSourceMissing
}

func LocalEnvPath() string {
	wd, err := os.Getwd()
	if err != nil {
		return ".env"
	}
	return filepath.Join(wd, ".env")
}

func readDotEnv(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	out := map[string]string{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		v = strings.Trim(v, "\"'")
		if k != "" {
			out[k] = v
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func UpsertKeyInEnvFile(path string, key string, value string) error {
	lines, err := readDotEnvLines(path)
	if err != nil && !errorsIsNotExist(err) {
		return err
	}

	updated := false
	for i, line := range lines {
		k, _, ok := parseEnvAssignment(line)
		if ok && k == key {
			lines[i] = fmt.Sprintf("%s=%s", key, value)
			updated = true
		}
	}
	if !updated {
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}
	return writeEnvLines(path, lines)
}

func RemoveKeyFromEnvFile(path string, key string) error {
	lines, err := readDotEnvLines(path)
	if err != nil {
		if errorsIsNotExist(err) {
			return nil
		}
		return err
	}

	out := make([]string, 0, len(lines))
	for _, line := range lines {
		k, _, ok := parseEnvAssignment(line)
		if ok && k == key {
			continue
		}
		out = append(out, line)
	}
	return writeEnvLines(path, out)
}

func readDotEnvLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := []string{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func writeEnvLines(path string, lines []string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	content := strings.Join(lines, "\n")
	if len(lines) > 0 {
		content += "\n"
	}
	return os.WriteFile(path, []byte(content), 0o600)
}

func parseEnvAssignment(line string) (string, string, bool) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return "", "", false
	}
	if strings.HasPrefix(trimmed, "export ") {
		trimmed = strings.TrimSpace(strings.TrimPrefix(trimmed, "export "))
	}
	parts := strings.SplitN(trimmed, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	k := strings.TrimSpace(parts[0])
	v := strings.Trim(strings.TrimSpace(parts[1]), "\"'")
	if k == "" {
		return "", "", false
	}
	return k, v, true
}

func errorsIsNotExist(err error) bool {
	return err != nil && os.IsNotExist(err)
}
