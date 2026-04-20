package timehelpers

import "time"

// Clock abstracts time so callers never call time.Now() directly.
type Clock interface {
	Now() time.Time
}

type realClock struct{}

// NewClock returns the production Clock backed by time.Now().
func NewClock() Clock { return &realClock{} }

func (realClock) Now() time.Time { return time.Now().Local() }
