package bag

import "errors"

// Snapshot de produto
type Item struct {
	id       string
	name     string
	imageKey string // Opcional
	quantity int32
}

func NewItem(id, name, imageKey string, quantity int32) (Item, error) {
	if id == "" {
		return Item{}, errors.New("id is required")
	}

	if name == "" {
		return Item{}, errors.New("name is required")
	}

	if quantity <= 0 {
		return Item{}, errors.New("quantity must be greater than 0")
	}

	return Item{
		id:       id,
		name:     name,
		imageKey: imageKey,
		quantity: quantity,
	}, nil
}

func (i Item) ID() string {
	return i.id
}

func (i Item) Name() string {
	return i.name
}

func (i Item) ImageKey() string {
	return i.imageKey
}

func (i Item) Quantity() int32 {
	return i.quantity
}
