package store

import "errors"

type Rating struct {
	sum   int32
	count int32
}

func NewRating(sum, count int32) Rating {
	return Rating{
		sum:   sum,
		count: count,
	}
}

func (r *Rating) Add(rating int16) error {
	if rating < 1 || rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	r.sum += int32(rating)
	r.count++

	return nil
}

func (r Rating) Average() float64 {
	if r.count == 0 {
		return 0
	}

	return float64(r.sum) / float64(r.count)
}

func (r Rating) Sum() int32 {
	return r.sum
}

func (r Rating) Count() int32 {
	return r.count
}
