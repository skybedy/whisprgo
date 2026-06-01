# START_NEW_AI_CHAT

Pokračujeme v projektu `whisprgo` v této složce:

`/home/skybedy/Programming/cli/golang/whisprgo`

Nejdřív načti kontext ze souborů:

- `README.md`
- `PROJECT_CONTEXT.md`
- `DECISIONS.md`
- `TODO.md`
- `AGENTS.md`
- `CHANGELOG.md`

Aktuální orientace:

- Primární cesta pro běžné diktování je raw transcript.
- Cleanup je volitelný (není hlavní cesta).
- Cleanup provider podporuje `openai`, `gemini`, `deepseek`.
- Secrets jsou provider/role based přes env + keyring (`whisprgo auth ...`).

Další cíl:

1. Navrhnout minimální implementaci lokální transkripce Parakeet.
2. Reuse OpenWhispr-style backend (`sherpa-onnx` websocket sidecar + ONNX model).
3. Nepřepisovat aplikaci od nuly; zachovat stávající pipeline `toggle -> transcribe -> output`.
