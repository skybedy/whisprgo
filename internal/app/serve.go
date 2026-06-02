package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"

	"whisprgo/internal/control"
	"whisprgo/internal/parakeet"
)

func (a *App) newServeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Run foreground local dictation backend",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runServe(cmd.Context(), cmd.ErrOrStderr())
		},
	}
}

func (a *App) runServe(parent context.Context, stderr interface{ Write([]byte) (int, error) }) error {
	ctx, stop := signal.NotifyContext(parent, os.Interrupt, syscall.SIGTERM)
	defer stop()

	a.foregroundLogs = true
	a.logInfof(stderr, "whisprgo serve starting")
	a.logInfof(stderr, "config loaded")

	if !stringsEqualFoldTrim(a.cfg.Transcription.Provider, "parakeet") {
		return fmt.Errorf("whisprgo serve currently requires transcription.provider=parakeet")
	}
	if a.parakeetMode() != "serve" {
		return fmt.Errorf("whisprgo serve requires transcription.parakeet.mode=serve")
	}

	cfg := a.cfg.Transcription.Parakeet
	a.logInfof(stderr, "transcription provider=parakeet")
	a.logInfof(stderr, "starting sherpa-onnx server")
	server, startupDuration, err := parakeet.StartManagedServer(ctx, cfg, parakeetAppLogger{app: a})
	if err != nil {
		return err
	}
	a.residentParakeet = server
	defer func() {
		a.logInfof(stderr, "stopping sherpa-onnx server")
		if stopErr := server.Stop(); stopErr != nil {
			a.logErrorf(stderr, "sherpa stop failed: %v", stopErr)
		}
		a.residentParakeet = nil
	}()
	a.logInfof(stderr, "sherpa server ready duration=%s", startupDuration)

	ln, err := control.ListenUnix()
	if err != nil {
		return err
	}
	defer func() {
		_ = ln.Close()
		_ = os.Remove(control.SocketPath())
	}()

	srv := &http.Server{Handler: a.controlMux(stderr)}
	go func() {
		<-ctx.Done()
		_ = srv.Shutdown(context.Background())
	}()
	go func() {
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			a.logErrorf(stderr, "control server failed: %v", err)
		}
	}()

	a.logInfof(stderr, "control socket listening at %s", control.SocketPath())
	a.logInfof(stderr, "waiting for commands")

	<-ctx.Done()
	if msg, err := a.performCancelLocal(stderr); err == nil && msg == "recording cancelled" {
		a.logInfof(stderr, "recording cancelled during shutdown")
	}
	return nil
}

func (a *App) controlMux(stderr interface{ Write([]byte) (int, error) }) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/toggle", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		noTranscribe := r.URL.Query().Get("no_transcribe") == "true"
		forceCleanup := r.URL.Query().Get("cleanup") == "true"
		res, err := a.performToggleLocal(r.Context(), stderr, noTranscribe, forceCleanup)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, control.Response{Message: res.Message, Audio: res.Audio, Text: res.Text})
	})
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		res, err := a.performStatusLocal(stderr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !res.Recording {
			writeJSON(w, http.StatusOK, control.Response{Message: "not recording"})
			return
		}
		writeJSON(w, http.StatusOK, control.Response{PID: res.PID, Audio: res.Audio, Started: res.Started})
	})
	mux.HandleFunc("/cancel", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		msg, err := a.performCancelLocal(stderr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, control.Response{Message: msg})
	})
	return mux
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func stringsEqualFoldTrim(a string, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}
