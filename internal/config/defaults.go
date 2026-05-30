package config

func Default() Config {
	return Config{
		Provider: ProviderConfig{
			Transcription:      "openai",
			TranscriptionModel: "gpt-4o-mini-transcribe",
		},
		Cleanup: CleanupConfig{
			Enabled:  false,
			Provider: "openai",
			Model:    "gpt-4o-mini",
			Prompt:   "Uprav tento diktovany text do prirozene cestiny.\nZachovej vyznam, nic nepridavej a nic duleziteho nemaz.",
		},
		Output: OutputConfig{
			Clipboard: true,
			Paste:     false,
		},
		Audio: AudioConfig{
			Input:    "default",
			Format:   "wav",
			Recorder: "ffmpeg",
		},
	}
}
