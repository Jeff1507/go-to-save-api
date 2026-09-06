package handler

type ItemResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	ImageURL *string `json:"imageUrl"`
	Quantity int32   `json:"quantity"`
}
