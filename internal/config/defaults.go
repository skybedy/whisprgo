package config

func Default() Config {
	return Config{
		Provider: "openai",
		Transcription: TranscriptionConfig{
			Model:          "whisper-1",
			Language:       "cs",
			TimeoutSeconds: 120,
		},
		Cleanup: CleanupConfig{
			Enabled:        true,
			Model:          "gpt-5-mini",
			TimeoutSeconds: 120,
			Prompt:         "Uprav tento diktovany text do prirozene cestiny. Zachovej vyznam.",
		},
		Audio: AudioConfig{
			InputDevice: "default",
			SampleRate:  16000,
			Channels:    1,
			Format:      "wav",
		},
		Output: OutputConfig{
			Mode:                "stdout",
			FilePath:            "",
			CopyToClipboard:     false,
			PasteToActiveWindow: false,
		},
		Security: SecurityConfig{
			SecretsBackend:   "keyring",
			AllowFileSecrets: false,
		},
	}
}
