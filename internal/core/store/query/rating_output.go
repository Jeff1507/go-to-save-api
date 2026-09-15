package query

import "math"

type RatingOutput struct {
	Average float64
	Count   int32
}

func calculateAverageRating(sum, count int32) float64 {
	if count == 0 {
		return 0
	}

	avg := float64(sum) / float64(count)

	return math.Round(avg*10) / 10
}
