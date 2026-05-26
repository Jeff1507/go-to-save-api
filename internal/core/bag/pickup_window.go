package bag

import (
	"errors"
	"time"
)

type PickupWindow struct {
	start time.Time
	end   time.Time
}

func NewPickupWindow(start, end time.Time) (PickupWindow, error) {
	const minInterval = time.Minute * 30 // 30 min
	const maxInterval = time.Hour * 6    // 06 hrs

	if start.IsZero() {
		return PickupWindow{}, errors.New("start is not set")
	}

	if end.IsZero() {
		return PickupWindow{}, errors.New("end is not set")
	}

	if end.Before(start) {
		return PickupWindow{}, errors.New("invalid pickup window")
	}

	if interval := end.Sub(start); interval < minInterval || interval > maxInterval {
		return PickupWindow{}, ErrInvalidPickupWindowInterval
	}

	return PickupWindow{
		start: start,
		end:   end,
	}, nil
}

func (pw PickupWindow) Start() time.Time {
	return pw.start
}

func (pw PickupWindow) End() time.Time {
	return pw.end
}
