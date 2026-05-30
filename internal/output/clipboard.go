package output

import (
	"fmt"
	"os/exec"
)

func CopyToClipboard(text string) error {
	if path, err := exec.LookPath("xclip"); err == nil {
		cmd := exec.Command(path, "-selection", "clipboard")
		in, err := cmd.StdinPipe()
		if err != nil {
			return fmt.Errorf("failed to open xclip stdin: %w", err)
		}
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start xclip: %w", err)
		}
		if _, err := in.Write([]byte(text)); err != nil {
			_ = in.Close()
			return fmt.Errorf("failed to write to xclip stdin: %w", err)
		}
		_ = in.Close()
		if err := cmd.Wait(); err != nil {
			return fmt.Errorf("xclip failed: %w", err)
		}
		return nil
	}

	if path, err := exec.LookPath("xsel"); err == nil {
		cmd := exec.Command(path, "--clipboard", "--input")
		in, err := cmd.StdinPipe()
		if err != nil {
			return fmt.Errorf("failed to open xsel stdin: %w", err)
		}
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start xsel: %w", err)
		}
		if _, err := in.Write([]byte(text)); err != nil {
			_ = in.Close()
			return fmt.Errorf("failed to write to xsel stdin: %w", err)
		}
		_ = in.Close()
		if err := cmd.Wait(); err != nil {
			return fmt.Errorf("xsel failed: %w", err)
		}
		return nil
	}

	return fmt.Errorf("clipboard is enabled but neither xclip nor xsel is installed")
}
