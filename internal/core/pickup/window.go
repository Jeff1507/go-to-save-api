package pickup

import (
	"errors"
	"time"
)

type Window struct {
	start time.Time
	end   time.Time
}

func NewWindow(start, end time.Time) (Window, error) {
	const minInterval = time.Minute * 30 // 30 min
	const maxInterval = time.Hour * 6    // 06 hrs

	if start.IsZero() {
		return Window{}, errors.New("start is not set")
	}

	if end.IsZero() {
		return Window{}, errors.New("end is not set")
	}

	if end.Before(start) {
		return Window{}, ErrWindowEndBeforeStart
	}

	if interval := end.Sub(start); interval < minInterval || interval > maxInterval {
		return Window{}, ErrInvalidWindowInterval
	}

	return Window{
		start: start,
		end:   end,
	}, nil
}

func (w Window) Extended(duration time.Duration) (Window, error) {
	return NewWindow(w.start, w.end.Add(duration))
}

func (w Window) IsZero() bool {
	return w.start.IsZero() && w.end.IsZero()
}

func (w Window) IsExpired(now time.Time) bool {
	return now.After(w.end)
}

func (w Window) Equal(other Window) bool {
	return w.start.Equal(other.start) && w.end.Equal(other.end)
}

func (w Window) IsNearExpiry(now time.Time, threshold time.Duration) bool {
	return w.end.Sub(now) <= threshold
}

func (w Window) Start() time.Time {
	return w.start
}

func (w Window) End() time.Time {
	return w.end
}
