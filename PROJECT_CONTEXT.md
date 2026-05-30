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
