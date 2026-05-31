package transcription

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const mistralTranscriptionURL = "https://api.mistral.ai/v1/audio/transcriptions"

type MistralProvider struct {
	apiKey string
	client *http.Client
}

func NewMistralProvider(apiKey string, client *http.Client) (*MistralProvider, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("MISTRAL_API_KEY is not set (env or keyring)")
	}
	if client == nil {
		client = &http.Client{}
	}
	return &MistralProvider{apiKey: strings.TrimSpace(apiKey), client: client}, nil
}

func (p *MistralProvider) Transcribe(ctx context.Context, audioPath string, model string) (string, error) {
	audioFile, err := os.Open(audioPath)
	if err != nil {
		return "", fmt.Errorf("failed to open audio file: %w", err)
	}
	defer audioFile.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err := writer.WriteField("model", model); err != nil {
		return "", fmt.Errorf("failed to write model field: %w", err)
	}

	fileWriter, err := writer.CreateFormFile("file", filepath.Base(audioPath))
	if err != nil {
		return "", fmt.Errorf("failed to create multipart file field: %w", err)
	}

	if _, err := io.Copy(fileWriter, audioFile); err != nil {
		return "", fmt.Errorf("failed to copy audio file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to finalize multipart body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, mistralTranscriptionURL, &body)
	if err != nil {
		return "", fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("transcription request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("transcription request failed with status %s: %s", resp.Status, string(raw))
	}

	text, err := extractTranscriptionText(raw)
	if err != nil {
		return "", err
	}
	if text == "" {
		return "", errors.New("empty transcription response")
	}

	return text, nil
}
