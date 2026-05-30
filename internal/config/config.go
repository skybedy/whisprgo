package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Provider ProviderConfig `yaml:"provider"`
	Cleanup  CleanupConfig  `yaml:"cleanup"`
	Output   OutputConfig   `yaml:"output"`
	Audio    AudioConfig    `yaml:"audio"`
}

type ProviderConfig struct {
	Transcription      string `yaml:"transcription"`
	TranscriptionModel string `yaml:"transcription_model"`
}

type CleanupConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Provider string `yaml:"provider"`
	Model    string `yaml:"model"`
	Prompt   string `yaml:"prompt"`
}

type OutputConfig struct {
	Clipboard      bool     `yaml:"clipboard"`
	Paste          bool     `yaml:"paste"`
	PasteDelayMs   int      `yaml:"paste_delay_ms"`
	PasteBlocklist []string `yaml:"paste_blocklist"`
}

type AudioConfig struct {
	Input    string `yaml:"input"`
	Format   string `yaml:"format"`
	Recorder string `yaml:"recorder"`
}

func Path() string {
	return "config.yaml"
}

func DisplayPath() string {
	return "./config.yaml"
}

func Load() (Config, error) {
	cfg := Default()
	path := Path()

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return Config{}, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
