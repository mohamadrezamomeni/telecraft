package state

import "time"

type State struct {
	Data       map[string]string
	Path       string
	Expiration time.Time
}
