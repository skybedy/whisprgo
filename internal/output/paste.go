package output

import (
	"fmt"
	"os/exec"
)

func PasteWithXdotool() error {
	path, err := exec.LookPath("xdotool")
	if err != nil {
		return fmt.Errorf("paste is enabled but xdotool is not installed")
	}

	cmd := exec.Command(path, "key", "--clearmodifiers", "ctrl+v")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("xdotool paste failed: %w", err)
	}

	return nil
}
