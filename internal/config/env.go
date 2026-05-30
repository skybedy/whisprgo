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

	paths := []string{
		LocalEnvPath(),
		EnvPath(),
	}
	for _, p := range paths {
		pairs, err := readDotEnv(p)
		if err != nil {
			continue
		}
		if v := strings.TrimSpace(pairs[key]); v != "" {
			return v
		}
	}

	return ""
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
