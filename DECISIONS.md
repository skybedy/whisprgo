# DECISIONS

- Project name is WhisprGo.
- Binary name is `whisprgo`.
- First target OS is Linux Mint Cinnamon / Ubuntu.
- No GUI in MVP.
- No Electron.
- No daemon in first version.
- No global hotkey handling inside the app.
- Keyboard shortcut will be configured externally in Linux Mint/Cinnamon.
- State file path is `/tmp/whisprgo/state.json`.
- Config file path is `<UserConfigDir>/whisprgo/config.yaml` (via `os.UserConfigDir`).
- API keys are never stored in config and are resolved from env first, then keyring.
- Recording will be delegated to external tools such as ffmpeg.
- Transcription and cleanup are separate concepts.
- Models are configurable through YAML config.
- Common settings are manageable through CLI config commands (`config path/show/get/set`).
- API keys are manageable through CLI auth commands (`auth status/set/delete`).
