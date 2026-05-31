# PROJECT_CONTEXT

- WhisperGo je minimalisticky Linux CLI nastroj v Go pro systemove diktovani.
- Cilove prostredi je Linux Mint Cinnamon / Ubuntu.
- Aplikace nema mit GUI, daemon ani vlastni globalni hotkeys v prvni verzi.
- Nahravani ma pozdeji ridit ffmpeg.
- Go aplikace bude pouze orchestrace: start recording, stop recording, transcribe, clipboard.
- Transcription model je audio -> text.
- Cleanup model je text -> text a je volitelny.
- Modely musi byt konfigurovatelne v YAML configu.
- API klice nepatri do configu.
- Primarni prvni workflow je `whispergo toggle`.
- Projekt ma pripraveny zaklad pro release: `VERSION`, prikaz `whispergo version`, smoke test a GitHub Actions CI/release workflow.
- Konfigurace je spravovatelna i pres CLI (`config init/show/get/set`), bez nutnosti rucni editace YAML.
- API key je spravovatelny pres CLI (`auth status/set-openai-key/clear-openai-key`) a uklada se do `~/.config/whispergo/.env`.
