package query

import "math"

type StoreOutput struct {
	ID        string
	Name      string
	AvatarURL *string
	Rating    RatingOutput
	Address   AddressOutput
}

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

type AddressOutput struct {
	PostalCode   string
	Street       string
	Number       string
	Neighborhood string
	Complement   *string
	City         string
	State        string
	Latitude     float64
	Longitude    float64
}
