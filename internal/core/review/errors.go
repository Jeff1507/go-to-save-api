package review

import "errors"

var (
	ErrAlreadyReviewed = errors.New("already reviewed")
	ErrCannotReview    = errors.New("can not review")
)
