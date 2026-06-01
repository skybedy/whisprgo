package recorder

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const recordsDir = "/tmp/whisprgo"

type FFmpegRecorder struct {
	opts Options
}

func NewFFmpegRecorder(opts Options) *FFmpegRecorder {
	if opts.InputDevice == "" {
		opts.InputDevice = "default"
	}
	if opts.SampleRate <= 0 {
		opts.SampleRate = 16000
	}
	if opts.Channels <= 0 {
		opts.Channels = 1
	}
	return &FFmpegRecorder{opts: opts}
}

func (r *FFmpegRecorder) CleanupOrphans() error {
	out, err := exec.Command("pgrep", "-f", "ffmpeg.*"+recordsDir+"/recording-.*\\.wav").Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil
		}
		return fmt.Errorf("failed to find orphan ffmpeg processes: %w", err)
	}

	for _, rawPID := range strings.Fields(string(out)) {
		pid, err := strconv.Atoi(rawPID)
		if err != nil || pid == os.Getpid() {
			continue
		}
		if !isFFmpegRecorderProcess(pid) {
			continue
		}
		proc, err := os.FindProcess(pid)
		if err != nil {
			continue
		}
		if err := proc.Signal(syscall.SIGTERM); err != nil && !errors.Is(err, os.ErrProcessDone) {
			return fmt.Errorf("failed to stop orphan ffmpeg process %d: %w", pid, err)
		}
	}
	return nil
}

func isFFmpegRecorderProcess(pid int) bool {
	raw, err := os.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "cmdline"))
	if err != nil {
		return false
	}
	cmdline := strings.ReplaceAll(string(raw), "\x00", " ")
	return strings.Contains(cmdline, "ffmpeg") &&
		strings.Contains(cmdline, recordsDir+"/recording-") &&
		strings.Contains(cmdline, ".wav")
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
	cmd := exec.Command(
		ffmpegPath,
		"-y",
		"-nostdin",
		"-loglevel", "error",
		"-f", "pulse",
		"-i", r.opts.InputDevice,
		"-ar", strconv.Itoa(r.opts.SampleRate),
		"-ac", strconv.Itoa(r.opts.Channels),
		"-c:a", "pcm_s16le",
		audioPath,
	)

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
	var (
		lastSize     int64 = -1
		stableRounds int
	)
	for {
		info, err := os.Stat(audioPath)
		if err == nil {
			size := info.Size()
			// WAV header is 44 bytes. We also require size stability to avoid
			// reading a partially flushed file right after SIGTERM.
			if size > 44 {
				if size == lastSize {
					stableRounds++
				} else {
					stableRounds = 0
					lastSize = size
				}
				if stableRounds >= 3 {
					return nil
				}
			}
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
