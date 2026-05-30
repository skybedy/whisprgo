package transcription

import "context"

type Provider interface {
	Transcribe(ctx context.Context, audioPath string, model string) (string, error)
}
