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

- [x] add VERSION and `whisprgo version`
- [x] add Linux release build targets (amd64/arm64)
- [x] add CHANGELOG and release checklist
- [x] add CLI smoke test script
- [x] add GitHub Actions CI and release workflows

## Config/Auth CLI management

- [x] add `config path/show/get/set`
- [x] validate supported config keys and value types
- [x] add `auth status/set/delete <provider>`
- [x] store secrets in keyring with env override fallback

## Security hardening (current)

- [x] move config path to OS user config dir
- [x] add env overrides (`WHISPRGO_*`) for selected non-secret keys
- [x] split secrets from config and use keyring + env priority
- [x] add `auth set/delete <provider>` for OpenAI/Gemini
- [x] block attempts to store secrets in `whisprgo config set ...`
- [x] add DeepSeek secret/provider wiring (`DEEPSEEK_API_KEY`, `auth set/delete deepseek`)

## Reliability/Latency hardening

- [x] unify log path to `~/.local/state/whisprgo/whisprgo.log`
- [x] parse multiple OpenAI response shapes for transcription and cleanup
- [x] fallback to raw transcription when cleanup fails
- [x] prevent overlapping `toggle` cycles while transcription is running
- [x] cleanup orphan ffmpeg recorders on new start
- [x] honor `audio.sample_rate` and `audio.channels` in ffmpeg recorder

## Next phase (Local transcription)

- [x] add local Parakeet transcription provider (no API key path)
- [x] reuse OpenWhispr-compatible `sherpa-onnx` websocket backend from Go
- [ ] add doctor checks for local Parakeet binary/model files
- [ ] keep cleanup optional; preserve fast raw-transcript path
