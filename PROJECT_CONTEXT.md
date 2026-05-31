# PROJECT_CONTEXT

- WhisprGo je minimalisticky Linux CLI nastroj v Go pro systemove diktovani.
- Cilove prostredi je Linux Mint Cinnamon / Ubuntu.
- Aplikace nema mit GUI, daemon ani vlastni globalni hotkeys v prvni verzi.
- Nahravani ma pozdeji ridit ffmpeg.
- Go aplikace bude pouze orchestrace: start recording, stop recording, transcribe, clipboard.
- Transcription model je audio -> text.
- Cleanup model je text -> text a je volitelny.
- Modely musi byt konfigurovatelne v YAML configu.
- API klice nepatri do configu.
- Primarni prvni workflow je `whisprgo toggle`.
- Projekt ma pripraveny zaklad pro release: `VERSION`, prikaz `whisprgo version`, smoke test a GitHub Actions CI/release workflow.
- Konfigurace je v OS user config dir (`$XDG_CONFIG_HOME/whisprgo/config.yaml`, fallback `~/.config/whisprgo/config.yaml`) a je spravovatelna pres CLI (`config path/show/get/set`).
- API keys jsou resene jako secrets: priorita je env (`OPENAI_API_KEY`, `GEMINI_API_KEY`) a fallback je systemovy keyring pres `auth status/set/delete`.
- Recorder respektuje audio parametry z configu (`audio.input_device`, `audio.sample_rate`, `audio.channels`), standardne 16kHz mono.
- Log path je stabilni: `~/.local/state/whisprgo/whisprgo.log`.
- Pri zpracovani predchozi nahravky se dalsi `toggle` ignoruje, aby nevznikaly prekrivajici se recording/transcribe behy.
