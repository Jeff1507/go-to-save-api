package handler

type RatingResponse struct {
	Average float64 `json:"average"`
	Count   int32   `json:"count"`
}
