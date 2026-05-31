package state

import "time"

type State struct {
	Recording      bool      `json:"recording"`
	Processing     bool      `json:"processing,omitempty"`
	PID            int       `json:"pid"`
	AudioPath      string    `json:"audio_path"`
	StartedAt      time.Time `json:"started_at"`
	NotificationID int       `json:"notification_id,omitempty"`
}
