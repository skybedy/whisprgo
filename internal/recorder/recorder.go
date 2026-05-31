package recorder

import "time"

type Session struct {
	PID       int
	AudioPath string
	StartedAt time.Time
}

type Options struct {
	InputDevice string
	SampleRate  int
	Channels    int
}

type Recorder interface {
	CleanupOrphans() error
	Start() (Session, error)
	Stop(pid int) error
	WaitForFile(audioPath string, timeout time.Duration) error
	IsRunning(pid int) bool
}
