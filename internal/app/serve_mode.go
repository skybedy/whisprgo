package app

import (
	"errors"
	"strings"
)

func (a *App) parakeetMode() string {
	mode := strings.ToLower(strings.TrimSpace(a.cfg.Transcription.Parakeet.Mode))
	if mode == "" {
		return "external"
	}
	return mode
}

func (a *App) shouldUseServeControl() bool {
	return strings.EqualFold(strings.TrimSpace(a.cfg.Transcription.Provider), "parakeet") && a.parakeetMode() == "serve"
}

func serveNotRunningError() error {
	return errors.New("WhisprGo serve nebezi. Spust nejdriv: whisprgo serve")
}
