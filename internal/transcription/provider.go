package transcription

import (
	"context"
	"errors"
)

// ErrEmptyTranscript is returned when the transcription backend returns no speech.
// It is not a failure — callers should handle it gracefully without error notifications.
var ErrEmptyTranscript = errors.New("no speech detected")

type Provider interface {
	Transcribe(ctx context.Context, audioPath string, model string) (string, error)
}
