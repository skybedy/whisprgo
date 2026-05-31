# WhisprGo

WhisprGo je minimalisticky Linux CLI nastroj v Go pro diktovani.
Nahraje zvuk z mikrofonu, posle ho na speech-to-text API a vysledek vlozi do schranky.

## Goals

- Vytvorit maly Linux CLI nastroj pro systemove diktovani.
- Mit jednoduchy workflow ovladany prikazem `whisprgo toggle`.
- Drzet MVP bez GUI, bez daemona a bez vlastnich hotkeys.

## Basic workflow

1. Prvni `whisprgo toggle` spusti nahravani.
2. Druhe `whisprgo toggle` nahravani zastavi.
3. Audio soubor se posle na transcription API.
4. Volitelne muze probehnout cleanup textu.
5. Vystup jde podle `output.mode` na stdout, do clipboardu, nebo se vlozi (`paste`) do aktivniho okna.

## Installation

```bash
make build
sudo make install
```

Update lokalni binarky:

```bash
make build
cp ./whisprgo ~/.local/bin/whisprgo
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

`$XDG_CONFIG_HOME/whisprgo/config.yaml`

Fallback:

`~/.config/whisprgo/config.yaml`

Pokud soubor neexistuje, aplikace pouzije defaulty v pameti.

## API key

API klic nepatri do configu.

Primarni varianta jsou environment variables:

- `OPENAI_API_KEY`
- `GEMINI_API_KEY`

Alternativa je systemovy keyring (`whisprgo auth ...`).

Priorita nacteni secretu:

1. env promenna
2. keyring
3. chyba

## Commands

- `whisprgo toggle`
- `whisprgo status`
- `whisprgo cancel`
- `whisprgo transcribe /path/to/audio.wav`
- `whisprgo config path`
- `whisprgo config show`
- `whisprgo config get <key>`
- `whisprgo config set <key> <value>`
- `whisprgo auth status`
- `whisprgo auth set openai`
- `whisprgo auth set gemini`
- `whisprgo auth delete openai`
- `whisprgo auth delete gemini`
- `whisprgo doctor`
- `whisprgo version`
- Global flag: `--verbose` (diagnosticke logy na stderr)

Podporovane config keys pro `config get/set`:

- `provider`
- `transcription.model`
- `transcription.language`
- `cleanup.enabled`
- `cleanup.model`
- `cleanup.prompt`
- `output.mode`
- `output.file_path`
- `output.copy_to_clipboard`
- `output.paste_to_active_window`

Minimalni priklad:

```yaml
provider: openai
transcription:
  model: whisper-1
  language: cs
cleanup:
  enabled: false
  model: gpt-5-mini
  prompt: Oprav pouze preklepy...
audio:
  input_device: default
  sample_rate: 16000
  channels: 1
output:
  mode: paste
security:
  secrets_backend: keyring
  allow_file_secrets: false
```

## Quick CLI setup (Linux Mint)

1. Diagnostika:

```bash
whisprgo doctor
```

2. Ulozeni API klice do keyringu:

```bash
whisprgo auth set openai
```

3. Doporucene rychle nastaveni:

```bash
whisprgo config set provider openai
whisprgo config set transcription.model whisper-1
whisprgo config set transcription.language cs
whisprgo config set cleanup.enabled false
whisprgo config set output.mode paste
```

4. Otestovani:

```bash
whisprgo toggle
whisprgo toggle
```

5. Volitelny cleanup (pomalejsi, ale kultivovanejsi text):

```bash
whisprgo config set cleanup.enabled true
whisprgo config set cleanup.model gpt-5-mini
whisprgo config set cleanup.prompt "Oprav pouze preklepy, interpunkci a zjevne chyby v diktovanem ceskem textu. Vrat pouze finalni opraveny text. Nevysvetluj, neptej se, nenabizej varianty, nepridavej odrazky ani zadne komentare. Pokud si nejsi jisty, zachovej puvodni formulaci."
```

## Development checks

```bash
make fmt
make test
go build ./cmd/whisprgo
scripts/smoke-test.sh
```

## Linux Mint keyboard shortcut

Aplikace sama neresi globalni hotkeys.
Hotkey se nastavi externe v Linux Mint/Cinnamon tak, aby volal:

`whisprgo toggle`

## State file

Stav nahravani je ulozen v:

`/tmp/whisprgo/state.json`

## Log file

Aplikace zapisuje provozni log vzdy do:

`~/.local/state/whisprgo/whisprgo.log`

## Recording indicator

Pri `toggle` start/stop aplikace posila desktop notifikaci pres `notify-send`, aby bylo jasne, zda nahravani bezi nebo uz bylo zastaveno.
Pri startu je notifikace sticky (zustane viditelna), pri stopu se nahradi na "Nahravani zastaveno.".
Pokud uz probiha transkripce predchozi nahravky, dalsi `toggle` nespusti novou nahravku.

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
