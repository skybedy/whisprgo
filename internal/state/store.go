package state

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	dirPath   = "/tmp/whisprgo"
	filePath  = "/tmp/whisprgo/state.json"
	fakeAudio = "/tmp/whisprgo/fake-recording.wav"
)

func Path() string {
	return filePath
}

func FakeAudioPath() string {
	return fakeAudio
}

func Exists() bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func Load() (State, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return State{}, err
	}

	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return State{}, err
	}

	return s, nil
}

func Save(s State) error {
	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0o644)
}

func Delete() error {
	err := os.Remove(filePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func EnsureDir() error {
	return os.MkdirAll(filepath.Dir(filePath), 0o755)
}
