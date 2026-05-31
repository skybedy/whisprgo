# WhisperGo

WhisperGo je minimalisticky Linux CLI nastroj v Go pro diktovani.
Nahraje zvuk z mikrofonu, posle ho na speech-to-text API a vysledek vlozi do schranky.

## Goals

- Vytvorit maly Linux CLI nastroj pro systemove diktovani.
- Mit jednoduchy workflow ovladany prikazem `whispergo toggle`.
- Drzet MVP bez GUI, bez daemona a bez vlastnich hotkeys.

## Basic workflow

1. Prvni `whispergo toggle` spusti nahravani.
2. Druhe `whispergo toggle` nahravani zastavi.
3. Audio soubor se posle na transcription API.
4. Vysledek jde do clipboardu.
5. Volitelne muze probehnout cleanup textu a paste pres `xdotool`.

## Installation

```bash
make build
sudo make install
```

Update lokalni binarky:

```bash
make build
cp ./whispergo ~/.local/bin/whispergo
```

Odinstalace z `/usr/local/bin`:

```bash
sudo make uninstall
```

Release build artifacts:

```bash
make release
ls -la dist/
```

## Requirements

- Linux (cilene Linux Mint Cinnamon / Ubuntu)
- Go 1.22+
- Pozdeji: `ffmpeg`, `xclip`/`xsel`, volitelne `xdotool`, `notify-send`

## Configuration

Konfigurace je v:

`./config.yaml`

Pokud soubor neexistuje, aplikace pouzije defaulty v pameti.

## API key

API klic nepatri do configu.

Primarni varianta je environment variable:

`OPENAI_API_KEY`

Volitelny fallback je soubor v projektu:

`./.env`

Aplikace umi i fallback:

`~/.config/whispergo/.env`

Priklad:

```bash
cp configs/.env.example ./.env
# potom uprav OPENAI_API_KEY v ./.env
```

Priorita nacteni:

1. `OPENAI_API_KEY` z prostredi
2. `OPENAI_API_KEY` z `./.env`
3. `OPENAI_API_KEY` z `~/.config/whispergo/.env`

## Commands

- `whispergo toggle`
- `whispergo status`
- `whispergo cancel`
- `whispergo transcribe /path/to/audio.wav`
- `whispergo config path`
- `whispergo config init [--force]`
- `whispergo config show`
- `whispergo config get <key>`
- `whispergo config set <key> <value>`
- `whispergo auth status`
- `whispergo auth set-openai-key`
- `whispergo auth clear-openai-key`
- `whispergo version`
- Global flag: `--verbose` (diagnosticke logy na stderr)

Podporovane config keys pro `config get/set`:

- `provider.transcription`
- `provider.transcription_model`
- `cleanup.enabled`
- `cleanup.provider`
- `cleanup.model`
- `cleanup.prompt`
- `output.clipboard`
- `output.paste`
- `output.paste_delay_ms`
- `audio.input`
- `audio.format`
- `audio.recorder`

Tip pro spolehlive paste pres hotkey:

```yaml
output:
  clipboard: true
  paste: true
  paste_delay_ms: 250
  paste_blocklist:
    - config.yaml
    - .env
```

Pokud je aktivni okno blokovane v `paste_blocklist`, text zustane v clipboardu, ale automaticky paste se neprovede.

## Development checks

```bash
make fmt
make test
go build ./cmd/whispergo
scripts/smoke-test.sh
```

## Linux Mint keyboard shortcut

Aplikace sama neresi globalni hotkeys.
Hotkey se nastavi externe v Linux Mint/Cinnamon tak, aby volal:

`whispergo toggle`

## State file

Stav nahravani je ulozen v:

`/tmp/whispergo/state.json`

## Log file

Aplikace zapisuje provozni log do:

`./whispergo.log`

Pokud nelze zapisovat do aktualni slozky, pouzije fallback:

`~/.local/state/whispergo/whispergo.log`

## Recording indicator

Pri `toggle` start/stop aplikace posila desktop notifikaci pres `notify-send`, aby bylo jasne, zda nahravani bezi nebo uz bylo zastaveno.
Pri startu je notifikace sticky (zustane viditelna), pri stopu se nahradi na "Nahravani zastaveno.".

## Project status

Aktualni faze: Phase 7.

- Phase 1: CLI skeleton, config loading, state handling.
- Phase 2: realne nahravani pres `ffmpeg` v `toggle/status/cancel`.
- Phase 3: `transcribe` vola OpenAI transcription API, model bere z configu a API klic z `OPENAI_API_KEY`.
- Phase 4: `toggle` po stopu automaticky vola transkripci (lze vypnout pres `--no-transcribe`).
- Phase 5: vystup lze kopirovat do clipboardu pres `xclip` nebo `xsel`.
- Phase 6: optional cleanup model (OpenAI) nad transkriptem.
- Phase 7: optional paste pres `xdotool`.

## Design principles

- Zadny GUI
- Zadny Electron
- Zadny daemon v prvni verzi
- Zadne hotkey handling primo v aplikaci
- Jednoducha a citelna Go implementace
