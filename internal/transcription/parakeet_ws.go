package transcription

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type ParakeetWSProvider struct {
	wsURL    string
	language string
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

	path := strings.TrimSpace(audioPath)
	if path == "" {
		return "", errors.New("audio path is empty")
	}
	payload, err := buildSherpaOfflinePayload(path)
	if err != nil {
		return "", err
	}
	if err := conn.WriteMessage(websocket.BinaryMessage, payload); err != nil {
		return "", fmt.Errorf("failed to send sherpa payload: %w", err)
	}
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return "", fmt.Errorf("failed to read sherpa response: %w", err)
	}
	text := extractSherpaText(msg)
	if text == "" {
		return "", ErrEmptyTranscript
	}
	_ = conn.WriteMessage(websocket.TextMessage, []byte("Done"))
	return text, nil
}

func extractSherpaText(msg []byte) string {
	type resp struct {
		Text string `json:"text"`
	}
	var r resp
	if err := json.Unmarshal(msg, &r); err == nil {
		return strings.TrimSpace(r.Text)
	}
	return strings.TrimSpace(string(msg))
}

func buildSherpaOfflinePayload(audioPath string) ([]byte, error) {
	var (
		samples    []float32
		sampleRate int
		err        error
	)
	for i := 0; i < 6; i++ {
		samples, sampleRate, err = readWavAsFloat32(audioPath)
		if err == nil {
			break
		}
		// ffmpeg can finish slightly after state switch; retry briefly for truncated wav tail.
		if strings.Contains(err.Error(), "invalid wav chunk size") || strings.Contains(err.Error(), "wav data chunk missing") {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	payload := make([]byte, 8+len(samples)*4)
	binary.LittleEndian.PutUint32(payload[0:4], uint32(sampleRate))
	binary.LittleEndian.PutUint32(payload[4:8], uint32(len(samples)*4))
	for i, v := range samples {
		binary.LittleEndian.PutUint32(payload[8+i*4:12+i*4], math.Float32bits(v))
	}
	return payload, nil
}

func readWavAsFloat32(path string) ([]float32, int, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read audio file: %w", err)
	}
	if len(raw) < 44 {
		return nil, 0, fmt.Errorf("unsupported wav file (too short): %s", filepath.Base(path))
	}
	if string(raw[0:4]) != "RIFF" || string(raw[8:12]) != "WAVE" {
		return nil, 0, errors.New("only WAVE files are supported for parakeet provider")
	}

	var (
		offset        = 12
		sampleRate    int
		numChannels   uint16
		bitsPerSample uint16
		audioFormat   uint16
		data          []byte
	)

	for offset+8 <= len(raw) {
		chunkID := string(raw[offset : offset+4])
		chunkSize := int(binary.LittleEndian.Uint32(raw[offset+4 : offset+8]))
		offset += 8
		if offset+chunkSize > len(raw) {
			if chunkID == "data" {
				chunkSize = len(raw) - offset
			} else {
				return nil, 0, errors.New("invalid wav chunk size")
			}
		}
		chunk := raw[offset : offset+chunkSize]
		switch chunkID {
		case "fmt ":
			if len(chunk) < 16 {
				return nil, 0, errors.New("invalid wav fmt chunk")
			}
			audioFormat = binary.LittleEndian.Uint16(chunk[0:2])
			numChannels = binary.LittleEndian.Uint16(chunk[2:4])
			sampleRate = int(binary.LittleEndian.Uint32(chunk[4:8]))
			bitsPerSample = binary.LittleEndian.Uint16(chunk[14:16])
		case "data":
			data = chunk
		}
		offset += chunkSize
		if chunkSize%2 == 1 {
			offset++
		}
	}

	if audioFormat != 1 {
		return nil, 0, errors.New("only PCM wav format is supported")
	}
	if numChannels != 1 {
		return nil, 0, errors.New("only mono wav is supported")
	}
	if bitsPerSample != 16 {
		return nil, 0, errors.New("only 16-bit wav is supported")
	}
	if len(data) == 0 {
		return nil, 0, errors.New("wav data chunk missing or invalid")
	}
	if len(data)%2 != 0 {
		data = data[:len(data)-1]
	}

	samples := make([]float32, len(data)/2)
	for i := 0; i < len(samples); i++ {
		s := int16(binary.LittleEndian.Uint16(data[i*2 : i*2+2]))
		samples[i] = float32(s) / 32768.0
	}
	return samples, sampleRate, nil
}
