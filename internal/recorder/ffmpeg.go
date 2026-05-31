package recorder

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

const recordsDir = "/tmp/whisprgo"

type FFmpegRecorder struct{}

func NewFFmpegRecorder() *FFmpegRecorder {
	return &FFmpegRecorder{}
}

func (r *FFmpegRecorder) Start() (Session, error) {
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return Session{}, errors.New("ffmpeg is not installed or not available in PATH")
	}

	if err := os.MkdirAll(recordsDir, 0o755); err != nil {
		return Session{}, fmt.Errorf("failed to create recording directory: %w", err)
	}

	audioPath := filepath.Join(recordsDir, fmt.Sprintf("recording-%s.wav", time.Now().Format("2006-01-02-150405.000000")))
	cmd := exec.Command(ffmpegPath, "-y", "-f", "pulse", "-i", "default", audioPath)

	if err := cmd.Start(); err != nil {
		return Session{}, fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	return Session{
		PID:       cmd.Process.Pid,
		AudioPath: audioPath,
		StartedAt: time.Now(),
	}, nil
}

func (r *FFmpegRecorder) Stop(pid int) error {
	if pid <= 0 {
		return fmt.Errorf("invalid pid: %d", pid)
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process %d: %w", pid, err)
	}

	if err := proc.Signal(syscall.SIGTERM); err != nil && !errors.Is(err, os.ErrProcessDone) {
		return fmt.Errorf("failed to stop process %d: %w", pid, err)
	}

	return nil
}

func (r *FFmpegRecorder) WaitForFile(audioPath string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		info, err := os.Stat(audioPath)
		if err == nil && info.Size() > 0 {
			return nil
		}
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to inspect audio file: %w", err)
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("audio file was not finalized in time: %s", audioPath)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (r *FFmpegRecorder) IsRunning(pid int) bool {
	if pid <= 0 {
		return false
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	err = proc.Signal(syscall.Signal(0))
	return err == nil
}
