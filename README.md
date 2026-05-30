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
4. Vysledek jde do clipboardu.
5. Volitelne muze probehnout cleanup textu a paste pres `xdotool`.

## Installation

```bash
go build -o whisprgo ./cmd/whisprgo
```

## Requirements

- Linux (cilene Linux Mint Cinnamon / Ubuntu)
- Go 1.22+
- Pozdeji: `ffmpeg`, `xclip`/`xsel`, volitelne `xdotool`

## Configuration

Konfigurace je v:

`~/.config/whisprgo/config.yaml`

Pokud soubor neexistuje, aplikace pouzije defaulty v pameti.

## API key

API klic nepatri do configu.

Pouzij environment variable:

`OPENAI_API_KEY`

## Commands

- `whisprgo toggle`
- `whisprgo status`
- `whisprgo cancel`
- `whisprgo transcribe /path/to/audio.wav`
- `whisprgo config path`

## Linux Mint keyboard shortcut

Aplikace sama neresi globalni hotkeys.
Hotkey se nastavi externe v Linux Mint/Cinnamon tak, aby volal:

`whisprgo toggle`

## State file

Stav nahravani je ulozen v:

`/tmp/whisprgo/state.json`

## Project status

Aktualni faze: Phase 1 (CLI skeleton, config loading, state handling, fake toggle).
Realne nahravani pres `ffmpeg`, API volani, clipboard a paste budou nasledovat v dalsich fazich.

## Design principles

- Zadny GUI
- Zadny Electron
- Zadny daemon v prvni verzi
- Zadne hotkey handling primo v aplikaci
- Jednoducha a citelna Go implementace
