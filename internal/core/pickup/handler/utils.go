package handler

import (
	"fmt"
	"time"
)

func parseTime(value string) (time.Time, error) {
	t, err := time.Parse("15:04", value)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse time: %v", err)
	}

	return t, nil
}
