# Changelog

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
