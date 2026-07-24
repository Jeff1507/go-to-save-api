package store

import "errors"

type Address struct {
	postalCode   string
	street       string
	number       string
	neighborhood string
	complement   string // Opcional
	city         string
	state        string
	latitude     float64
	longitude    float64
}

func NewAddress(
	postalCode string,
	street string,
	number string,
	neighborhood string,
	complement string,
	city string,
	state string,
	latitude float64,
	longitude float64,
) (Address, error) {
	if postalCode == "" {
		return Address{}, errors.New("zip is required")
	}

	if street == "" {
		return Address{}, errors.New("street is required")
	}

	if number == "" {
		return Address{}, errors.New("number is required")
	}

	if neighborhood == "" {
		return Address{}, errors.New("neighborhood is required")
	}

	if city == "" {
		return Address{}, errors.New("city is required")
	}

	if state == "" {
		return Address{}, errors.New("state is required")
	}

	if latitude < -90 || latitude > 90 {
		return Address{}, errors.New("latitude out of range")
	}

	if longitude < -180 || longitude > 180 {
		return Address{}, errors.New("longitude out of range")
	}

	return Address{
		postalCode:   postalCode,
		street:       street,
		number:       number,
		neighborhood: neighborhood,
		complement:   complement,
		city:         city,
		state:        state,
		latitude:     latitude,
		longitude:    longitude,
	}, nil
}

func (a Address) PostalCode() string {
	return a.postalCode
}

func (a Address) Street() string {
	return a.street
}

func (a Address) Number() string {
	return a.number
}

func (a Address) Neighborhood() string {
	return a.neighborhood
}

func (a Address) Complement() string {
	return a.complement
}

func (a Address) City() string {
	return a.city
}

func (a Address) State() string {
	return a.state
}

func (a Address) Latitude() float64 {
	return a.latitude
}

func (a Address) Longitude() float64 {
	return a.longitude
}
