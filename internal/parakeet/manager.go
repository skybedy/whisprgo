package parakeet

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"whisprgo/internal/config"
)

type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
}

type nopLogger struct{}

func (nopLogger) Infof(string, ...any)  {}
func (nopLogger) Errorf(string, ...any) {}

type ManagedServer struct {
	cmd    *exec.Cmd
	host   string
	port   int
	logger Logger
}

func StartManagedServer(ctx context.Context, cfg config.ParakeetConfig, logger Logger) (*ManagedServer, time.Duration, error) {
	if logger == nil {
		logger = nopLogger{}
	}
	if err := ValidateManagedConfig(cfg); err != nil {
		return nil, 0, err
	}
	addr := net.JoinHostPort(strings.TrimSpace(cfg.Host), fmt.Sprintf("%d", cfg.Port))
	if err := ensurePortAvailable(addr); err != nil {
		return nil, 0, err
	}

	cmd := exec.CommandContext(
		ctx,
		strings.TrimSpace(cfg.Binary),
		fmt.Sprintf("--tokens=%s", filepath.Join(cfg.ModelDir, "tokens.txt")),
		fmt.Sprintf("--encoder=%s", filepath.Join(cfg.ModelDir, "encoder.int8.onnx")),
		fmt.Sprintf("--decoder=%s", filepath.Join(cfg.ModelDir, "decoder.int8.onnx")),
		fmt.Sprintf("--joiner=%s", filepath.Join(cfg.ModelDir, "joiner.int8.onnx")),
		fmt.Sprintf("--port=%d", cfg.Port),
		fmt.Sprintf("--num-threads=%d", cfg.NumThreads),
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to capture parakeet stdout: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to capture parakeet stderr: %w", err)
	}

	startedAt := time.Now()
	if err := cmd.Start(); err != nil {
		return nil, 0, fmt.Errorf("failed to start parakeet server: %w", err)
	}

	server := &ManagedServer{
		cmd:    cmd,
		host:   strings.TrimSpace(cfg.Host),
		port:   cfg.Port,
		logger: logger,
	}

	go streamPipe("parakeet stdout", stdout, logger.Infof)
	go streamPipe("parakeet stderr", stderr, logger.Infof)

	if err := waitForTCPReady(addr, time.Duration(cfg.StartupTimeoutSeconds)*time.Second); err != nil {
		_ = server.Stop()
		return nil, 0, err
	}

	return server, time.Since(startedAt), nil
}

func (s *ManagedServer) URL() string {
	return fmt.Sprintf("ws://%s:%d", s.host, s.port)
}

func (s *ManagedServer) Stop() error {
	if s == nil || s.cmd == nil || s.cmd.Process == nil {
		return nil
	}
	if err := s.cmd.Process.Signal(syscall.SIGTERM); err != nil && !errors.Is(err, os.ErrProcessDone) {
		_ = s.cmd.Process.Kill()
		return fmt.Errorf("failed to stop parakeet server: %w", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- s.cmd.Wait()
	}()

	select {
	case err := <-done:
		if err != nil && !isExpectedExit(err) {
			return fmt.Errorf("parakeet server exited with error: %w", err)
		}
		return nil
	case <-time.After(2 * time.Second):
		_ = s.cmd.Process.Kill()
		err := <-done
		if err != nil && !isExpectedExit(err) {
			return fmt.Errorf("parakeet server did not stop cleanly: %w", err)
		}
		return nil
	}
}

func ValidateManagedConfig(cfg config.ParakeetConfig) error {
	mode := strings.ToLower(strings.TrimSpace(cfg.Mode))
	if mode != "managed" {
		return fmt.Errorf("unsupported parakeet mode for managed validation: %s", cfg.Mode)
	}
	if strings.TrimSpace(cfg.Binary) == "" {
		return errors.New("parakeet managed binary is empty")
	}
	if stat, err := os.Stat(cfg.Binary); err != nil || stat.IsDir() {
		return fmt.Errorf("parakeet binary not found: %s", cfg.Binary)
	}
	if strings.TrimSpace(cfg.ModelDir) == "" {
		return errors.New("parakeet managed model_dir is empty")
	}
	if stat, err := os.Stat(cfg.ModelDir); err != nil || !stat.IsDir() {
		return fmt.Errorf("parakeet model_dir not found: %s", cfg.ModelDir)
	}
	for _, name := range []string{"tokens.txt", "encoder.int8.onnx", "decoder.int8.onnx", "joiner.int8.onnx"} {
		path := filepath.Join(cfg.ModelDir, name)
		if stat, err := os.Stat(path); err != nil || stat.IsDir() {
			return fmt.Errorf("missing parakeet model file: %s", path)
		}
	}
	if strings.TrimSpace(cfg.Host) == "" {
		return errors.New("parakeet host is empty")
	}
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("parakeet port is invalid: %d", cfg.Port)
	}
	if cfg.NumThreads <= 0 {
		return fmt.Errorf("parakeet num_threads is invalid: %d", cfg.NumThreads)
	}
	if cfg.StartupTimeoutSeconds <= 0 {
		return fmt.Errorf("parakeet startup_timeout_seconds is invalid: %d", cfg.StartupTimeoutSeconds)
	}
	if cfg.RequestTimeoutSeconds <= 0 {
		return fmt.Errorf("parakeet request_timeout_seconds is invalid: %d", cfg.RequestTimeoutSeconds)
	}
	return nil
}

func CheckExternalEndpoint(wsURL string, timeout time.Duration) error {
	u, err := url.Parse(strings.TrimSpace(wsURL))
	if err != nil {
		return fmt.Errorf("invalid sherpa_ws_url: %w", err)
	}
	if u.Scheme != "ws" && u.Scheme != "wss" {
		return fmt.Errorf("invalid sherpa_ws_url scheme: %s", u.Scheme)
	}
	addr := u.Host
	if !strings.Contains(addr, ":") {
		switch u.Scheme {
		case "ws":
			addr = net.JoinHostPort(addr, "80")
		case "wss":
			addr = net.JoinHostPort(addr, "443")
		}
	}
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return fmt.Errorf("sherpa endpoint is not reachable: %w", err)
	}
	_ = conn.Close()
	return nil
}

func waitForTCPReady(addr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		conn, err := net.DialTimeout("tcp", addr, 250*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("Parakeet managed server did not start within %s. Check binary path, model_dir and port.", timeout)
		}
		time.Sleep(150 * time.Millisecond)
	}
}

func ensurePortAvailable(addr string) error {
	conn, err := net.DialTimeout("tcp", addr, 250*time.Millisecond)
	if err == nil {
		_ = conn.Close()
		return fmt.Errorf("parakeet port is already in use: %s", addr)
	}
	return nil
}

func streamPipe(prefix string, r io.Reader, logf func(format string, args ...any)) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		logf("%s: %s", prefix, line)
	}
}

func isExpectedExit(err error) bool {
	if err == nil {
		return true
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		return false
	}
	return exitErr.ExitCode() == -1 || exitErr.ExitCode() == 0 || exitErr.ExitCode() == 143
}
