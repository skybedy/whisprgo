package transcription

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

type ParakeetWSProvider struct {
	wsURL    string
	language string
}

type parakeetWSRequest struct {
	Type     string `json:"type,omitempty"`
	AudioPath string `json:"audio_path"`
	Model    string `json:"model,omitempty"`
	Language string `json:"language,omitempty"`
}

type parakeetWSResponse struct {
	Type  string `json:"type,omitempty"`
	Text  string `json:"text,omitempty"`
	Error string `json:"error,omitempty"`
}

func NewParakeetWSProvider(wsURL, language string) (*ParakeetWSProvider, error) {
	u := strings.TrimSpace(wsURL)
	if u == "" {
		return nil, errors.New("transcription.sherpa_ws_url is empty")
	}
	if !strings.HasPrefix(u, "ws://") && !strings.HasPrefix(u, "wss://") {
		return nil, errors.New("transcription.sherpa_ws_url must start with ws:// or wss://")
	}
	return &ParakeetWSProvider{wsURL: u, language: strings.TrimSpace(language)}, nil
}

func (p *ParakeetWSProvider) Transcribe(ctx context.Context, audioPath string, model string) (string, error) {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.DialContext(ctx, p.wsURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to connect to sherpa websocket: %w", err)
	}
	defer conn.Close()

	req := parakeetWSRequest{
		Type:      "transcribe",
		AudioPath: strings.TrimSpace(audioPath),
		Model:     strings.TrimSpace(model),
		Language:  p.language,
	}
	if req.AudioPath == "" {
		return "", errors.New("audio path is empty")
	}
	if err := conn.WriteJSON(req); err != nil {
		return "", fmt.Errorf("failed to send sherpa request: %w", err)
	}

	for {
		var raw json.RawMessage
		if err := conn.ReadJSON(&raw); err != nil {
			return "", fmt.Errorf("failed to read sherpa response: %w", err)
		}
		text, done, err := parseParakeetWSResponse(raw)
		if err != nil {
			return "", err
		}
		if done {
			return text, nil
		}
	}
}

func parseParakeetWSResponse(raw json.RawMessage) (text string, done bool, err error) {
	var resp parakeetWSResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return "", false, fmt.Errorf("failed to parse sherpa response: %w", err)
	}
	if strings.TrimSpace(resp.Error) != "" {
		return "", false, fmt.Errorf("sherpa transcription failed: %s", strings.TrimSpace(resp.Error))
	}
	if t := strings.TrimSpace(resp.Text); t != "" {
		return t, true, nil
	}
	if strings.EqualFold(strings.TrimSpace(resp.Type), "result") {
		return "", false, errors.New("sherpa response has empty text")
	}
	return "", false, nil
}
