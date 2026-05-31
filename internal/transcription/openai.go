package transcription

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"whisprgo/internal/config"
)

const openAITranscriptionURL = "https://api.openai.com/v1/audio/transcriptions"

type OpenAIProvider struct {
	apiKey string
	client *http.Client
}

func NewOpenAIProviderFromEnv(client *http.Client) (*OpenAIProvider, error) {
	apiKey := config.ResolveSecret("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY is not set (env or ~/.config/whisprgo/.env)")
	}

	if client == nil {
		client = &http.Client{}
	}

	return &OpenAIProvider{apiKey: apiKey, client: client}, nil
}

func (p *OpenAIProvider) Transcribe(ctx context.Context, audioPath string, model string) (string, error) {
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, openAITranscriptionURL, &body)
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

	var parsed struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("failed to parse transcription response: %w", err)
	}

	if parsed.Text == "" {
		return "", errors.New("empty transcription response")
	}

	return parsed.Text, nil
}
