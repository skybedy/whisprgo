package transcription

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
)

func TestReadWavAsFloat32(t *testing.T) {
	dir := t.TempDir()
	wav := filepath.Join(dir, "test.wav")
	writeTestWavMono16(t, wav, 16000, []int16{0, 16384, -16384})

	samples, sr, err := readWavAsFloat32(wav)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sr != 16000 {
		t.Fatalf("unexpected sample rate: %d", sr)
	}
	if len(samples) != 3 {
		t.Fatalf("unexpected sample count: %d", len(samples))
	}
}

func TestBuildSherpaOfflinePayload(t *testing.T) {
	dir := t.TempDir()
	wav := filepath.Join(dir, "test.wav")
	writeTestWavMono16(t, wav, 8000, []int16{0, 1000})

	payload, err := buildSherpaOfflinePayload(wav)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := binary.LittleEndian.Uint32(payload[0:4]); got != 8000 {
		t.Fatalf("unexpected sample rate header: %d", got)
	}
	if got := binary.LittleEndian.Uint32(payload[4:8]); got != 8 {
		t.Fatalf("unexpected payload byte length: %d", got)
	}
}

func TestExtractSherpaText(t *testing.T) {
	got := extractSherpaText([]byte(`{"text":" ahoj svete "}`))
	if got != "ahoj svete" {
		t.Fatalf("unexpected text: %q", got)
	}
	got = extractSherpaText([]byte("plain text"))
	if got != "plain text" {
		t.Fatalf("unexpected fallback text: %q", got)
	}
}

func writeTestWavMono16(t *testing.T, path string, sampleRate int, samples []int16) {
	t.Helper()
	dataSize := len(samples) * 2
	riffSize := 36 + dataSize
	buf := make([]byte, 44+dataSize)
	copy(buf[0:4], []byte("RIFF"))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(riffSize))
	copy(buf[8:12], []byte("WAVE"))
	copy(buf[12:16], []byte("fmt "))
	binary.LittleEndian.PutUint32(buf[16:20], 16)
	binary.LittleEndian.PutUint16(buf[20:22], 1)
	binary.LittleEndian.PutUint16(buf[22:24], 1)
	binary.LittleEndian.PutUint32(buf[24:28], uint32(sampleRate))
	binary.LittleEndian.PutUint32(buf[28:32], uint32(sampleRate*2))
	binary.LittleEndian.PutUint16(buf[32:34], 2)
	binary.LittleEndian.PutUint16(buf[34:36], 16)
	copy(buf[36:40], []byte("data"))
	binary.LittleEndian.PutUint32(buf[40:44], uint32(dataSize))
	o := 44
	for _, s := range samples {
		binary.LittleEndian.PutUint16(buf[o:o+2], uint16(s))
		o += 2
	}
	if err := os.WriteFile(path, buf, 0o644); err != nil {
		t.Fatalf("failed to write wav: %v", err)
	}
}
