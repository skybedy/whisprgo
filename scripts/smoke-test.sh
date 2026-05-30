#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

cleanup() {
  go run ./cmd/whisprgo cancel >/dev/null 2>&1 || true
}
trap cleanup EXIT

# Basic commands
go run ./cmd/whisprgo version >/dev/null
go run ./cmd/whisprgo config path >/dev/null
go run ./cmd/whisprgo status >/dev/null

# Toggle lifecycle without API dependency
go run ./cmd/whisprgo toggle --no-transcribe >/dev/null
go run ./cmd/whisprgo status | grep -q "recording"
go run ./cmd/whisprgo toggle --no-transcribe >/dev/null
go run ./cmd/whisprgo status | grep -q "not recording"

# Transcribe validation error path (missing file)
set +e
out=$(go run ./cmd/whisprgo transcribe /tmp/whisprgo-missing-file.wav 2>&1)
code=$?
set -e
if [[ $code -eq 0 ]]; then
  echo "expected transcribe missing-file to fail"
  exit 1
fi
if ! grep -q "audio file does not exist" <<<"$out"; then
  echo "expected missing file error, got:"
  echo "$out"
  exit 1
fi

echo "smoke test passed"
