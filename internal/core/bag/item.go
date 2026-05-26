package bag

import "errors"

type Item struct {
	id       string
	quantity int32
}

func NewItem(id string, quantity int32) (Item, error) {
	if id == "" {
		return Item{}, errors.New("id is required")
	}

	if quantity <= 0 {
		return Item{}, errors.New("quantity must be greater than 0")
	}

	return Item{
		id:       id,
		quantity: quantity,
	}, nil
}

func (i Item) ID() string {
	return i.id
}

func (i Item) Quantity() int32 {
	return i.quantity
}
