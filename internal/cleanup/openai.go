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

	"whisprgo/internal/secrets"
)

const openAIResponsesURL = "https://api.openai.com/v1/responses"

type OpenAICleaner struct {
	apiKey string
	client *http.Client
}

func NewOpenAICleaner(apiKey string, client *http.Client) (*OpenAICleaner, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("OPENAI_API_KEY is not set (env or keyring)")
	}
	if client == nil {
		client = &http.Client{}
	}
	return &OpenAICleaner{apiKey: strings.TrimSpace(apiKey), client: client}, nil
}

func NewOpenAICleanerFromSecrets(client *http.Client) (*OpenAICleaner, error) {
	apiKey := func() string {
		v, _, err := secrets.Get("openai")
		if err != nil {
			return ""
		}
		return v
	}()
	return NewOpenAICleaner(apiKey, client)
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

	text, err := extractCleanupText(rawResp)
	if err != nil {
		return "", err
	}
	if text == "" {
		return "", errors.New("empty cleanup response")
	}

	return text, nil
}

func extractCleanupText(raw []byte) (string, error) {
	var byOutputText struct {
		OutputText string `json:"output_text"`
	}
	if err := json.Unmarshal(raw, &byOutputText); err == nil {
		if v := strings.TrimSpace(byOutputText.OutputText); v != "" {
			return v, nil
		}
	}

	var byText struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(raw, &byText); err == nil {
		if v := strings.TrimSpace(byText.Text); v != "" {
			return v, nil
		}
	}

	var generic map[string]any
	if err := json.Unmarshal(raw, &generic); err != nil {
		return "", fmt.Errorf("failed to parse cleanup response: %w", err)
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
				if txt, _ := cm["text"].(string); strings.TrimSpace(txt) != "" {
					parts = append(parts, strings.TrimSpace(txt))
				}
			}
		}
		if len(parts) > 0 {
			return strings.Join(parts, "\n"), nil
		}
	}

	return "", nil
}
