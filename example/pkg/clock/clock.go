package clock

import "time"

//nolint: lll
//go:generate mockgen -destination mock/clock_mock.go github.com/018bf/example/pkg/clock Clock

type Clock interface {
	Now() time.Time
}

type RealClock struct {
}

func NewRealClock() Clock {
	return &RealClock{}
}

func (c RealClock) Now() time.Time {
	return time.Now()
}
