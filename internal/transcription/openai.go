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
	"strings"

	"whisprgo/internal/secrets"
)

const openAITranscriptionURL = "https://api.openai.com/v1/audio/transcriptions"

type OpenAIProvider struct {
	apiKey string
	client *http.Client
}

func NewOpenAIProvider(apiKey string, client *http.Client) (*OpenAIProvider, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("OPENAI_API_KEY is not set (env or keyring)")
	}
	if client == nil {
		client = &http.Client{}
	}
	return &OpenAIProvider{apiKey: strings.TrimSpace(apiKey), client: client}, nil
}

func NewOpenAIProviderFromSecrets(client *http.Client) (*OpenAIProvider, error) {
	apiKey := func() string {
		v, _, err := secrets.Get("openai")
		if err != nil {
			return ""
		}
		return v
	}()
	return NewOpenAIProvider(apiKey, client)
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

	text, err := extractTranscriptionText(raw)
	if err != nil {
		return "", err
	}
	if text == "" {
		return "", errors.New("empty transcription response")
	}

	return text, nil
}

func extractTranscriptionText(raw []byte) (string, error) {
	var byText struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(raw, &byText); err == nil {
		if v := strings.TrimSpace(byText.Text); v != "" {
			return v, nil
		}
	}

	var byOutputText struct {
		OutputText string `json:"output_text"`
	}
	if err := json.Unmarshal(raw, &byOutputText); err == nil {
		if v := strings.TrimSpace(byOutputText.OutputText); v != "" {
			return v, nil
		}
	}

	var generic map[string]any
	if err := json.Unmarshal(raw, &generic); err != nil {
		return "", fmt.Errorf("failed to parse transcription response: %w", err)
	}

	if output, ok := generic["output"].([]any); ok {
		var parts []string
		for _, item := range output {
			obj, ok := item.(map[string]any)
			if !ok {
				continue
			}
			content, ok := obj["content"].([]any)
			if !ok {
				continue
			}
			for _, c := range content {
				cm, ok := c.(map[string]any)
				if !ok {
					continue
				}
				if t, _ := cm["type"].(string); t == "output_text" {
					if txt, _ := cm["text"].(string); strings.TrimSpace(txt) != "" {
						parts = append(parts, strings.TrimSpace(txt))
					}
				}
			}
		}
		if len(parts) > 0 {
			return strings.Join(parts, "\n"), nil
		}
	}

	return "", nil
}
