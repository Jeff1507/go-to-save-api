package order

import "errors"

type Item struct {
	bagID     string
	unitPrice int64 // Preço em centavos
	quantity  int32
}

func NewItem(bagID string, unitPrice int64, quantity int32) (Item, error) {
	if bagID == "" {
		return Item{}, errors.New("bag id is required")
	}

	if unitPrice <= 0 {
		return Item{}, errors.New("unit price must be greater than 0")
	}

	if quantity <= 0 {
		return Item{}, errors.New("quantity must be greater than 0")
	}

	return Item{
		bagID:     bagID,
		unitPrice: unitPrice,
		quantity:  quantity,
	}, nil
}

func (i Item) TotalPrice() int64 {
	return i.unitPrice * int64(i.quantity)
}

func (i Item) BagID() string {
	return i.bagID
}

func (i Item) UnitPrice() int64 {
	return i.unitPrice
}

func (i Item) Quantity() int32 {
	return i.quantity
}
