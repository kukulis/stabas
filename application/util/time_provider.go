package util

import "time"

type TimeProvider interface {
	ProvideTime() time.Time
}

type SimpleTimeProvider struct{}

func (sp SimpleTimeProvider) ProvideTime() time.Time {
	return time.Now()
}

type FixedTimeProvider struct {
	Time time.Time
}

func (sp FixedTimeProvider) ProvideTime() time.Time {
	return sp.Time
}
