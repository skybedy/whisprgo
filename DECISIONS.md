# DECISIONS

- Project name is WhisperGo.
- Binary name is `whispergo`.
- First target OS is Linux Mint Cinnamon / Ubuntu.
- No GUI in MVP.
- No Electron.
- No daemon in first version.
- No global hotkey handling inside the app.
- Keyboard shortcut will be configured externally in Linux Mint/Cinnamon.
- State file path is `/tmp/whispergo/state.json`.
- Config file path is `./config.yaml`.
- API keys are stored in environment variables, not in config.
- Recording will be delegated to external tools such as ffmpeg.
- Transcription and cleanup are separate concepts.
- Models are configurable through YAML config.
- Common settings are manageable through CLI config commands (`config init/show/get/set`).
- API keys are manageable through CLI auth commands and are stored outside `config.yaml`.
