package pickup

import (
	"errors"
	"time"
)

type Time struct {
	hour   int
	minute int
}

func NewTime(hour, minute int) (Time, error) {
	if hour < 0 || hour >= 24 {
		return Time{}, errors.New("hour out of range")
	}

	if minute < 0 || minute >= 60 {
		return Time{}, errors.New("minute out of range")
	}

	return Time{
		hour:   hour,
		minute: minute,
	}, nil
}

func (t Time) Duration() time.Duration {
	return time.Duration(t.hour)*time.Hour + time.Duration(t.minute)*time.Minute
}

func (t Time) String() string {
	return FormatTime(t.hour, t.minute)
}

func (t Time) Hour() int {
	return t.hour
}

func (t Time) Minute() int {
	return t.minute
}
