package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func EnvPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "~/.config/whisprgo/.env"
	}
	return filepath.Join(home, ".config", "whisprgo", ".env")
}

func ResolveSecret(key string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	pairs, err := readDotEnv(EnvPath())
	if err != nil {
		return ""
	}
	return strings.TrimSpace(pairs[key])
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
