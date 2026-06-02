package config

func Default() Config {
	return Config{
		Provider: "openai",
		Transcription: TranscriptionConfig{
			Provider:    "openai",
			Model:       "whisper-1",
			Language:    "cs",
			SherpaWSURL: "ws://127.0.0.1:6006",
			Parakeet: ParakeetConfig{
				Mode:                  "external",
				SherpaWSURL:           "ws://127.0.0.1:6006",
				Host:                  "127.0.0.1",
				Port:                  6010,
				NumThreads:            4,
				StartupTimeoutSeconds: 30,
				RequestTimeoutSeconds: 120,
			},
			TimeoutSeconds: 120,
		},
		Cleanup: CleanupConfig{
			Enabled:        false,
			Provider:       "openai",
			Model:          "gpt-5-mini",
			TimeoutSeconds: 120,
			Prompt:         "Oprav pouze preklepy, interpunkci a zjevne chyby v diktovanem ceskem textu. Vrat pouze finalni opraveny text. Nevysvetluj, neptej se, nenabizej varianty, nepridavej odrazky ani zadne komentare. Pokud si nejsi jisty, zachovej puvodni formulaci.",
		},
		Audio: AudioConfig{
			InputDevice:        "default",
			SampleRate:         16000,
			Channels:           1,
			Format:             "wav",
			MaxDurationSeconds: 900,
		},
		Output: OutputConfig{
			Mode:                "stdout",
			FilePath:            "",
			CopyToClipboard:     false,
			PasteToActiveWindow: false,
		},
		Logging: LoggingConfig{
			IncludeText: false,
			FilePath:    "",
		},
		Security: SecurityConfig{
			SecretsBackend:   "keyring",
			AllowFileSecrets: false,
		},
	}
}
