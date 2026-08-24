package handler

type AddressResponse struct {
	PostalCode   string  `json:"postalCode"`
	Street       string  `json:"street"`
	Number       string  `json:"number"`
	Neighborhood string  `json:"neighborhood"`
	Complement   *string `json:"complement"`
	City         string  `json:"city"`
	State        string  `json:"state"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
}
