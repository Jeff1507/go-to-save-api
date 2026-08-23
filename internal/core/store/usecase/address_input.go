package usecase

type AddressInput struct {
	PostalCode   string
	Street       string
	Number       string
	Neighborhood string
	Complement   string
	City         string
	State        string
	Latitude     float64
	Longitude    float64
}
