# Changelog

## Unreleased

### Added
- DeepSeek cleanup provider support (`cleanup.provider=deepseek`).
- DeepSeek secret wiring via env/keyring (`DEEPSEEK_API_KEY`, `auth set/delete deepseek`).

### Changed
- Config path moved to OS user config dir (`$XDG_CONFIG_HOME/whisprgo/config.yaml`, fallback `~/.config/whisprgo/config.yaml`).
- Auth flow migrated to keyring-backed secrets with env priority (`OPENAI_API_KEY`, `GEMINI_API_KEY`).
- CLI auth commands updated to `auth status`, `auth set <provider>`, `auth delete <provider>`.
- Runtime log path unified to `~/.local/state/whisprgo/whisprgo.log`.
- Recorder now respects configured audio capture params (`audio.sample_rate`, `audio.channels`) to reduce upload size and latency.

### Fixed
- Transcription and cleanup response parsing now supports multiple OpenAI response shapes (`text`, `output_text`, `output[].content[].text`).
- Cleanup failure no longer blocks final output; app falls back to raw transcription.
- `toggle` no longer starts overlapping recordings while previous transcription is in progress.
- Recorder now cleans up orphan ffmpeg capture processes on new start.

## 0.1.0 - 2026-05-30

### Added
- Complete MVP phases 1-7.
- Real ffmpeg recording lifecycle (`toggle`, `status`, `cancel`).
- OpenAI transcription command with model from YAML config.
- Optional cleanup model via OpenAI Responses API.
- Optional clipboard output (`xclip`/`xsel`) and optional paste (`xdotool`).
- Verbose diagnostics and file logging (`~/.local/state/whisprgo/whisprgo.log`).
- `version` command and build metadata (version/commit/date).
- Makefile with build/install/test/release targets.
- Baseline tests for config loading and state store lifecycle.
