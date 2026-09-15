package handler

type BagItemResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	ImageURL *string `json:"imageUrl"`
	Quantity int32   `json:"quantity"`
}

type BagResponse struct {
	ID        string            `json:"id"`
	Title     string            `json:"title"`
	Quantity  int32             `json:"quantity"`
	UnitPrice int64             `json:"unitPrice"`
	Items     []BagItemResponse `json:"items"`
}
