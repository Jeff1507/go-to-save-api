package handler

type StoreResponse struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	AvatarURL *string         `json:"avatarUrl"`
	Rating    RatingResponse  `json:"rating"`
	Address   AddressResponse `json:"address"`
}

type RatingResponse struct {
	Average float64 `json:"average"`
	Count   int32   `json:"count"`
}

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
