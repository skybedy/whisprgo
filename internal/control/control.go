package control

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Response struct {
	Message string `json:"message,omitempty"`
	Text    string `json:"text,omitempty"`
	PID     int    `json:"pid,omitempty"`
	Audio   string `json:"audio,omitempty"`
	Started string `json:"started,omitempty"`
}

func SocketPath() string {
	runtimeDir := strings.TrimSpace(os.Getenv("XDG_RUNTIME_DIR"))
	if runtimeDir != "" {
		return filepath.Join(runtimeDir, "whisprgo.sock")
	}
	return filepath.Join(os.TempDir(), "whisprgo", "whisprgo.sock")
}

func EnsureSocketDir() error {
	return os.MkdirAll(filepath.Dir(SocketPath()), 0o755)
}

func ListenUnix() (net.Listener, error) {
	if err := EnsureSocketDir(); err != nil {
		return nil, err
	}
	if err := removeStaleSocket(); err != nil {
		return nil, err
	}
	return net.Listen("unix", SocketPath())
}

func Client(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				var d net.Dialer
				return d.DialContext(ctx, "unix", SocketPath())
			},
		},
	}
}

func IsServeReachable(ctx context.Context) bool {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://unix/health", nil)
	resp, err := Client(time.Second).Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func SendToggle(ctx context.Context, noTranscribe bool, forceCleanup bool) (Response, error) {
	path := "/toggle"
	params := make([]string, 0, 2)
	if noTranscribe {
		params = append(params, "no_transcribe=true")
	}
	if forceCleanup {
		params = append(params, "cleanup=true")
	}
	if len(params) > 0 {
		path += "?" + strings.Join(params, "&")
	}
	return doJSONRequest(ctx, http.MethodPost, path, nil)
}

func SendCancel(ctx context.Context) (Response, error) {
	return doJSONRequest(ctx, http.MethodPost, "/cancel", nil)
}

func SendStatus(ctx context.Context) (Response, error) {
	return doJSONRequest(ctx, http.MethodGet, "/status", nil)
}

func doJSONRequest(ctx context.Context, method string, path string, body any) (Response, error) {
	var payload io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return Response{}, err
		}
		payload = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(ctx, method, "http://unix"+path, payload)
	if err != nil {
		return Response{}, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := Client(10 * time.Minute).Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := strings.TrimSpace(string(raw))
		if msg == "" {
			msg = resp.Status
		}
		return Response{}, fmt.Errorf(msg)
	}
	var out Response
	if len(raw) == 0 {
		return out, nil
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return Response{}, err
	}
	return out, nil
}

func removeStaleSocket() error {
	path := SocketPath()
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if info.Mode()&os.ModeSocket == 0 {
		return fmt.Errorf("control path exists and is not a socket: %s", path)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	if IsServeReachable(ctx) {
		return fmt.Errorf("WhisprGo serve already runs at %s", path)
	}
	return os.Remove(path)
}
