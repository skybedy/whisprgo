package cleanup

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const geminiGenerateContentURL = "https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s"

type GeminiCleaner struct {
	apiKey string
	client *http.Client
}

func NewGeminiCleaner(apiKey string, client *http.Client) (*GeminiCleaner, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("GEMINI_API_KEY is not set (env or keyring)")
	}
	if client == nil {
		client = &http.Client{}
	}
	return &GeminiCleaner{apiKey: strings.TrimSpace(apiKey), client: client}, nil
}

func (c *GeminiCleaner) Clean(ctx context.Context, input string, model string, prompt string) (string, error) {
	if strings.TrimSpace(model) == "" {
		return "", errors.New("cleanup.model is empty")
	}
	endpoint := fmt.Sprintf(geminiGenerateContentURL, url.PathEscape(model), url.QueryEscape(c.apiKey))

	payload := map[string]any{
		"contents": []map[string]any{
			{
				"role": "user",
				"parts": []map[string]string{
					{"text": prompt + "\n\n" + input},
				},
			},
		},
	}

	rawBody, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to encode cleanup request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(rawBody))
	if err != nil {
		return "", fmt.Errorf("failed to build cleanup request: %w", err)
	}
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

	text, err := extractGeminiCleanupText(rawResp)
	if err != nil {
		return "", err
	}
	if text == "" {
		return "", errors.New("empty cleanup response")
	}
	return text, nil
}

func extractGeminiCleanupText(raw []byte) (string, error) {
	var resp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return "", fmt.Errorf("failed to parse cleanup response: %w", err)
	}

	var parts []string
	for _, c := range resp.Candidates {
		for _, p := range c.Content.Parts {
			if t := strings.TrimSpace(p.Text); t != "" {
				parts = append(parts, t)
			}
		}
	}
	return strings.TrimSpace(strings.Join(parts, "\n")), nil
}
