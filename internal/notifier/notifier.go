package notifier

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Notifier struct {
	notifySendPath string
}

func New() *Notifier {
	path, err := exec.LookPath("notify-send")
	if err != nil {
		return &Notifier{}
	}
	return &Notifier{notifySendPath: path}
}

func (n *Notifier) Available() bool {
	return n.notifySendPath != ""
}

func (n *Notifier) Notify(summary string, body string, icon string) error {
	_, err := n.NotifyWithOptions(summary, body, icon, -1, 0, false)
	return err
}

func (n *Notifier) NotifyWithOptions(summary string, body string, icon string, expireMs int, replaceID int, wantID bool) (int, error) {
	if n.notifySendPath == "" {
		return 0, nil
	}

	args := []string{"-a", "WhisperGo"}
	if wantID {
		args = append(args, "-p")
	}
	if replaceID > 0 {
		args = append(args, "-r", strconv.Itoa(replaceID))
	}
	if expireMs >= 0 {
		args = append(args, "-t", strconv.Itoa(expireMs))
	}
	if icon != "" {
		args = append(args, "-i", icon)
	}
	args = append(args, summary, body)

	cmd := exec.Command(n.notifySendPath, args...)
	if wantID {
		out, err := cmd.Output()
		if err != nil {
			return 0, fmt.Errorf("notify-send failed: %w", err)
		}
		idStr := strings.TrimSpace(string(out))
		if idStr == "" {
			return 0, nil
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return 0, nil
		}
		return id, nil
	}
	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("notify-send failed: %w", err)
	}
	return 0, nil
}
