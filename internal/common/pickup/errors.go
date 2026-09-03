package pickup

import "errors"

var (
	ErrInvalidWindowInterval = errors.New("invalid pickup window interval")
	ErrWindowEndBeforeStart  = errors.New("pickup window end must be after start")
	ErrWindowStartInPast     = errors.New("pickup window must not start in the past")
)
