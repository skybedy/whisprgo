# TODO

## Phase 1

- [x] CLI skeleton
- [x] YAML config
- [x] state handling
- [x] fake toggle

## Phase 2

- [x] real ffmpeg recording
- [x] status with running PID
- [x] cancel real recording

## Phase 3

- [x] transcribe existing audio file through OpenAI
- [x] read model from config
- [x] read API key from OPENAI_API_KEY

## Phase 4

- [x] connect toggle stop to transcription
- [x] print result to stdout

## Phase 5

- [x] clipboard output via xclip/xsel

## Phase 6

- [x] optional cleanup model

## Phase 7

- [x] optional paste via xdotool

## Post-MVP (Release readiness)

- [x] add VERSION and `whispergo version`
- [x] add Linux release build targets (amd64/arm64)
- [x] add CHANGELOG and release checklist
- [x] add CLI smoke test script
- [x] add GitHub Actions CI and release workflows

## Config/Auth CLI management

- [x] add `config init/show/get/set`
- [x] validate supported config keys and value types
- [x] add `auth status/set-openai-key/clear-openai-key`
- [x] store OpenAI key in `~/.config/whispergo/.env` with secure permissions
