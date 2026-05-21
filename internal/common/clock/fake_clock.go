package clock

import "time"

type FakeClock struct {
	DateTime time.Time
}

func NewFakeClock() FakeClock {
	return FakeClock{
		DateTime: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
}

func (fc FakeClock) Now() time.Time {
	return fc.DateTime
}
