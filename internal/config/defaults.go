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
			Prompt:         "Oprav pouze preklepy, interpunkci a zjevne chyby v diktovanem ceskem textu. Vrat pouze finalni opraveny text. Nevysvetluj, neptej se, nenabizej varianty, nepridavej odrazky ani zadne komentare. Pokud si nejsi jisty, zachovej puvodni formulaci.",
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
