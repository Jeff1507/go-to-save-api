package review

import (
	"errors"
	"time"
)

type Review struct {
	id         string
	orderID    string
	rating     int16
	comment    string // Opcional
	reviewedAt time.Time
}

func New(id, orderID string, rating int16, comment string, now time.Time) (*Review, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	if orderID == "" {
		return nil, errors.New("order id is required")
	}

	if rating < 1 || rating > 5 {
		return nil, errors.New("rating must be between 1 and 5")
	}

	if now.IsZero() {
		return nil, errors.New("now is not set")
	}

	return &Review{
		id:         id,
		orderID:    orderID,
		rating:     rating,
		comment:    comment,
		reviewedAt: now,
	}, nil
}

func From(id, orderID string, rating int16, comment string, reviewedAt time.Time) *Review {
	return &Review{
		id:         id,
		orderID:    orderID,
		rating:     rating,
		comment:    comment,
		reviewedAt: reviewedAt,
	}
}

func (r Review) ID() string {
	return r.id
}

func (r Review) OrderID() string {
	return r.orderID
}

func (r Review) Rating() int16 {
	return r.rating
}

func (r Review) Comment() string {
	return r.comment
}

func (r Review) ReviewedAt() time.Time {
	return r.reviewedAt
}
