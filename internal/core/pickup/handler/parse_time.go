package handler

import "time"

func parseTime(value string) time.Time {
	t, _ := time.Parse("15:04", value)
	return t
}
