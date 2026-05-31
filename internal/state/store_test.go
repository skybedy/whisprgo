package state

import (
	"os"
	"testing"
	"time"
)

func preserveStateFile(t *testing.T) {
	t.Helper()

	original, readErr := os.ReadFile(Path())
	hadOriginal := readErr == nil

	t.Cleanup(func() {
		_ = Delete()
		if hadOriginal {
			_ = os.MkdirAll("/tmp/whisprgo", 0o755)
			_ = os.WriteFile(Path(), original, 0o644)
		}
	})
}

func TestSaveLoadDeleteLifecycle(t *testing.T) {
	preserveStateFile(t)
	_ = Delete()

	started := time.Now().UTC().Truncate(time.Second)
	in := State{
		Recording: true,
		PID:       4242,
		AudioPath: "/tmp/whisprgo/test.wav",
		StartedAt: started,
	}

	if err := Save(in); err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}
	if !Exists() {
		t.Fatalf("Exists() should be true after Save()")
	}

	out, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if out.PID != in.PID || out.AudioPath != in.AudioPath || !out.StartedAt.Equal(in.StartedAt) {
		t.Fatalf("Load() mismatch: got %+v want %+v", out, in)
	}

	if err := Delete(); err != nil {
		t.Fatalf("Delete() returned error: %v", err)
	}
	if Exists() {
		t.Fatalf("Exists() should be false after Delete()")
	}
}
