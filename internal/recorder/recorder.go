package recorder

import "time"

type Session struct {
	PID       int
	AudioPath string
	StartedAt time.Time
}

type Recorder interface {
	Start() (Session, error)
	Stop(pid int) error
	IsRunning(pid int) bool
}
