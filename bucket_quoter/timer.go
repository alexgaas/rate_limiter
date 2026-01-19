package bucket_quoter

import "time"

// Timers - only ms now (add us by demand)

type InstantTimer interface {
	Now() int64
	Duration(from int64, to int64) int64
	Resolution() int64
}

type InstantTimerMs struct {
	resolution int64 // milliseconds
}

func NewInstantTimerMs() *InstantTimerMs {
	return &InstantTimerMs{
		resolution: 1000,
	}
}

func (t *InstantTimerMs) Now() int64 {
	return time.Now().UnixMilli()
}

func (t *InstantTimerMs) Duration(from int64, to int64) int64 {
	return to - from
}

func (t *InstantTimerMs) Resolution() int64 {
	return t.resolution
}
