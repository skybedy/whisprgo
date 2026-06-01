package cleanup

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const deepSeekChatCompletionsURL = "https://api.deepseek.com/chat/completions"

type DeepSeekCleaner struct {
	apiKey string
	client *http.Client
}

func NewDeepSeekCleaner(apiKey string, client *http.Client) (*DeepSeekCleaner, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("DEEPSEEK_API_KEY is not set (env or keyring)")
	}
	if client == nil {
		client = &http.Client{}
	}
	return &DeepSeekCleaner{apiKey: strings.TrimSpace(apiKey), client: client}, nil
}

func (c *DeepSeekCleaner) Clean(ctx context.Context, input string, model string, prompt string) (string, error) {
	payload := map[string]any{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": prompt},
			{"role": "user", "content": input},
		},
	}

	rawBody, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to encode cleanup request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, deepSeekChatCompletionsURL, bytes.NewReader(rawBody))
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

	text, err := extractDeepSeekCleanupText(rawResp)
	if err != nil {
		return "", err
	}
	if text == "" {
		return "", errors.New("empty cleanup response")
	}
	return text, nil
}

func extractDeepSeekCleanupText(raw []byte) (string, error) {
	var resp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return "", fmt.Errorf("failed to parse cleanup response: %w", err)
	}
	if len(resp.Choices) == 0 {
		return "", nil
	}
	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}
