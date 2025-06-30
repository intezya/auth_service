package clock

import "time"

type Clock interface {
	Now() time.Time
	Unix(sec int64, nsec int64) time.Time
}

type RealClock struct{}

func NewRealClock() Clock {
	return &RealClock{}
}

func (c *RealClock) Now() time.Time {
	return time.Now()
}

func (c *RealClock) Unix(sec int64, nsec int64) time.Time {
	return time.Unix(sec, nsec)
}

type MockClock struct {
	currentTime time.Time
}

func NewMockClock(t time.Time) *MockClock {
	return &MockClock{currentTime: t}
}

func (c *MockClock) Now() time.Time {
	return c.currentTime
}

func (c *MockClock) Unix(sec int64, nsec int64) time.Time {
	return time.Unix(sec, nsec)
}

func (c *MockClock) SetTime(t time.Time) {
	c.currentTime = t
}
