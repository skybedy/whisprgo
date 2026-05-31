package cleanup

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"whisprgo/internal/secrets"
)

const openAIResponsesURL = "https://api.openai.com/v1/responses"

type OpenAICleaner struct {
	apiKey string
	client *http.Client
}

func NewOpenAICleanerFromSecrets(client *http.Client) (*OpenAICleaner, error) {
	apiKey := func() string {
		v, _, err := secrets.Get("openai")
		if err != nil {
			return ""
		}
		return v
	}()
	if apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY is not set (env or keyring)")
	}
	if client == nil {
		client = &http.Client{}
	}
	return &OpenAICleaner{apiKey: apiKey, client: client}, nil
}

func (c *OpenAICleaner) Clean(ctx context.Context, input string, model string, prompt string) (string, error) {
	payload := map[string]any{
		"model": model,
		"input": []map[string]any{
			{
				"role": "system",
				"content": []map[string]string{
					{"type": "input_text", "text": prompt},
				},
			},
			{
				"role": "user",
				"content": []map[string]string{
					{"type": "input_text", "text": input},
				},
			},
		},
	}

	rawBody, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to encode cleanup request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIResponsesURL, bytes.NewReader(rawBody))
	if err != nil {
		return "", fmt.Errorf("failed to build cleanup request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("cleanup request failed: %w", err)
	}
	defer resp.Body.Close()

	rawResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read cleanup response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("cleanup request failed with status %s: %s", resp.Status, string(rawResp))
	}

	var parsed struct {
		OutputText string `json:"output_text"`
	}
	if err := json.Unmarshal(rawResp, &parsed); err != nil {
		return "", fmt.Errorf("failed to parse cleanup response: %w", err)
	}
	if parsed.OutputText == "" {
		return "", errors.New("empty cleanup response")
	}

	return parsed.OutputText, nil
}
