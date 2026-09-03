package pickup

import (
	"errors"
	"time"
)

var brazilLocation = time.FixedZone("BRT", -3*60*60)

func ResolveWindow(start, end Time, now time.Time) (Window, error) {
	if now.IsZero() {
		return Window{}, errors.New("now is not set")
	}

	today := now.In(brazilLocation)
	baseDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, brazilLocation)

	startAt := baseDate.Add(start.Duration())
	endAt := baseDate.Add(end.Duration())

	if startAt.Before(now) {
		return Window{}, ErrWindowStartInPast
	}

	if end.Duration() <= start.Duration() {
		endAt = endAt.Add(time.Hour * 24)
	}

	return NewWindow(startAt, endAt)
}
