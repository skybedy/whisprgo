package output

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func PasteWithXdotool(delay time.Duration, blocklist []string) (bool, string, error) {
	path, err := exec.LookPath("xdotool")
	if err != nil {
		return false, "", fmt.Errorf("paste is enabled but xdotool is not installed")
	}

	if delay > 0 {
		time.Sleep(delay)
	}

	title, err := activeWindowTitle(path)
	if err != nil {
		return false, "", err
	}
	for _, blocked := range blocklist {
		blocked = strings.TrimSpace(blocked)
		if blocked != "" && strings.Contains(strings.ToLower(title), strings.ToLower(blocked)) {
			return false, title, nil
		}
	}

	cmd := exec.Command(path, "key", "--clearmodifiers", "ctrl+v")
	if err := cmd.Run(); err != nil {
		return false, title, fmt.Errorf("xdotool paste failed: %w", err)
	}

	return true, title, nil
}

func activeWindowTitle(xdotoolPath string) (string, error) {
	out, err := exec.Command(xdotoolPath, "getactivewindow", "getwindowname").Output()
	if err != nil {
		return "", fmt.Errorf("failed to read active window title: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
