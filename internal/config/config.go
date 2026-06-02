package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Provider      string              `yaml:"provider"`
	Transcription TranscriptionConfig `yaml:"transcription"`
	Cleanup       CleanupConfig       `yaml:"cleanup"`
	Audio         AudioConfig         `yaml:"audio"`
	Output        OutputConfig        `yaml:"output"`
	Logging       LoggingConfig       `yaml:"logging"`
	Security      SecurityConfig      `yaml:"security"`
}

type TranscriptionConfig struct {
	Provider       string         `yaml:"provider"`
	Model          string         `yaml:"model"`
	Language       string         `yaml:"language"`
	SherpaWSURL    string         `yaml:"sherpa_ws_url"`
	Parakeet       ParakeetConfig `yaml:"parakeet"`
	TimeoutSeconds int            `yaml:"timeout_seconds"`
}

type ParakeetConfig struct {
	Mode                  string `yaml:"mode"`
	SherpaWSURL           string `yaml:"sherpa_ws_url"`
	Binary                string `yaml:"binary"`
	ModelDir              string `yaml:"model_dir"`
	Host                  string `yaml:"host"`
	Port                  int    `yaml:"port"`
	NumThreads            int    `yaml:"num_threads"`
	StartupTimeoutSeconds int    `yaml:"startup_timeout_seconds"`
	RequestTimeoutSeconds int    `yaml:"request_timeout_seconds"`
}

type CleanupConfig struct {
	Enabled        bool   `yaml:"enabled"`
	Provider       string `yaml:"provider"`
	Model          string `yaml:"model"`
	TimeoutSeconds int    `yaml:"timeout_seconds"`
	Prompt         string `yaml:"prompt"`
}

type AudioConfig struct {
	InputDevice        string `yaml:"input_device"`
	SampleRate         int    `yaml:"sample_rate"`
	Channels           int    `yaml:"channels"`
	Format             string `yaml:"format"`
	MaxDurationSeconds int    `yaml:"max_duration_seconds"`
}

type OutputConfig struct {
	Mode                string `yaml:"mode"`
	FilePath            string `yaml:"file_path"`
	CopyToClipboard     bool   `yaml:"copy_to_clipboard"`
	PasteToActiveWindow bool   `yaml:"paste_to_active_window"`
}

type SecurityConfig struct {
	SecretsBackend   string `yaml:"secrets_backend"`
	AllowFileSecrets bool   `yaml:"allow_file_secrets"`
}

type LoggingConfig struct {
	IncludeText bool   `yaml:"include_text"`
	FilePath    string `yaml:"file_path"`
}

func Path() string {
	cfgDir, err := os.UserConfigDir()
	if err != nil || cfgDir == "" {
		return filepath.Join(".", "config.yaml")
	}
	return filepath.Join(cfgDir, "whisprgo", "config.yaml")
}

func DisplayPath() string {
	return Path()
}

func Load() (Config, error) {
	cfg := Default()
	path := Path()

	data, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return Config{}, err
		}
	} else {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return Config{}, err
		}
	}

	ApplyEnvOverrides(&cfg)
	normalizeProviders(&cfg)
	return cfg, nil
}

func normalizeProviders(cfg *Config) {
	if cfg == nil {
		return
	}
	global := cfg.Provider
	if global == "" {
		global = "openai"
	}
	if cfg.Transcription.Provider == "" {
		cfg.Transcription.Provider = global
	}
	if cfg.Cleanup.Provider == "" {
		cfg.Cleanup.Provider = global
	}
	if cfg.Transcription.Parakeet.Mode == "" {
		cfg.Transcription.Parakeet.Mode = "external"
	}
	if cfg.Transcription.SherpaWSURL != "" && (cfg.Transcription.Parakeet.SherpaWSURL == "" || cfg.Transcription.Parakeet.SherpaWSURL == Default().Transcription.Parakeet.SherpaWSURL) {
		cfg.Transcription.Parakeet.SherpaWSURL = cfg.Transcription.SherpaWSURL
	}
	if cfg.Transcription.SherpaWSURL == "" && cfg.Transcription.Parakeet.SherpaWSURL != "" {
		cfg.Transcription.SherpaWSURL = cfg.Transcription.Parakeet.SherpaWSURL
	}
	if cfg.Transcription.Parakeet.Host == "" {
		cfg.Transcription.Parakeet.Host = "127.0.0.1"
	}
	if cfg.Transcription.Parakeet.Port <= 0 {
		cfg.Transcription.Parakeet.Port = 6010
	}
	if cfg.Transcription.Parakeet.NumThreads <= 0 {
		cfg.Transcription.Parakeet.NumThreads = 4
	}
	if cfg.Transcription.Parakeet.StartupTimeoutSeconds <= 0 {
		cfg.Transcription.Parakeet.StartupTimeoutSeconds = 30
	}
	if cfg.Transcription.Parakeet.RequestTimeoutSeconds <= 0 {
		cfg.Transcription.Parakeet.RequestTimeoutSeconds = 120
	}
}

func Exists() bool {
	_, err := os.Stat(Path())
	return err == nil
}

func Save(cfg Config) error {
	raw, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}
	path := Path()
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, raw, 0o644)
}
