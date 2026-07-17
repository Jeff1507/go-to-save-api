package bag

import (
	"errors"
	"fmt"
	"time"
)

type Status uint8

const (
	Unknown Status = iota
	Active
	Closed
)

func (s Status) String() string {
	switch s {
	case Active:
		return "ACTIVE"
	case Closed:
		return "CLOSED"
	default:
		return "UNKNOWN"
	}
}

func ParseStatus(str string) (Status, error) {
	switch str {
	case "ACTIVE":
		return Active, nil
	case "CLOSED":
		return Closed, nil
	default:
		return Unknown, fmt.Errorf("unknown status: %q", str)
	}
}

type Bag struct {
	id            string
	storeID       string
	title         string
	description   string // Opcional
	stock         int32
	originalPrice int64 // Preço em centavos
	finalPrice    int64 // Preço em centavos
	status        Status
	pickupWindow  PickupWindow
	items         []Item
	createdAt     time.Time
	updatedAt     time.Time
}

func calculateFinalPrice(originalPrice int64, discount int) int64 {
	return originalPrice - (originalPrice * int64(discount) / 100)
}

func unique(items []Item) []Item {
	seen := make(map[string]bool)
	lst := make([]Item, 0, len(items))

	for _, item := range items {
		if !seen[item.name] {
			seen[item.name] = true
			lst = append(lst, item)
		}
	}

	return lst
}

func New(
	id string,
	storeID string,
	title string,
	description string,
	stock int32,
	originalPrice int64,
	discount int,
	pickupWindow PickupWindow,
	items []Item,
	now time.Time,
) (*Bag, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	if storeID == "" {
		return nil, errors.New("store id is required")
	}

	if title == "" {
		return nil, errors.New("title is required")
	}

	if stock <= 0 {
		return nil, errors.New("stock must be greater than 0")
	}

	if originalPrice <= 0 {
		return nil, errors.New("original price must be greater than 0")
	}

	if discount < 30 || discount > 90 {
		return nil, errors.New("discount must be between 30% and 90%")
	}

	if len(items) <= 0 {
		return nil, errors.New("at least one item is required")
	}

	if now.IsZero() {
		return nil, errors.New("now is not set")
	}

	return &Bag{
		id:            id,
		storeID:       storeID,
		title:         title,
		description:   description,
		stock:         stock,
		originalPrice: originalPrice,
		finalPrice:    calculateFinalPrice(originalPrice, discount),
		status:        Active,
		pickupWindow:  pickupWindow,
		items:         unique(items),
		createdAt:     now,
		updatedAt:     now,
	}, nil
}

func From(
	id string,
	storeID string,
	title string,
	description string,
	stock int32,
	originalPrice int64,
	finalPrice int64,
	status Status,
	pickupWindow PickupWindow,
	items []Item,
	createdAt time.Time,
	updatedAt time.Time,
) *Bag {
	return &Bag{
		id:            id,
		storeID:       storeID,
		title:         title,
		description:   description,
		stock:         stock,
		originalPrice: originalPrice,
		finalPrice:    finalPrice,
		status:        status,
		pickupWindow:  pickupWindow,
		items:         items,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}

func (b Bag) HasStock() bool {
	return b.stock > 0
}

func (b Bag) IsExpired(now time.Time) bool {
	return now.After(b.pickupWindow.end)
}

func (b Bag) IsClosed() bool {
	return b.status == Closed
}

func (b *Bag) Close(now time.Time) error {
	if now.IsZero() {
		return errors.New("now is not set")
	}

	if b.IsClosed() {
		return ErrInvalidOperation
	}

	b.status = Closed
	b.updatedAt = now

	return nil
}

func (b *Bag) ChangeTitle(title string, now time.Time) error {
	if title == "" {
		return errors.New("title is required")
	}

	if now.IsZero() {
		return errors.New("now is not set")
	}

	b.title = title
	b.updatedAt = now

	return nil
}

func (b *Bag) ChangeDescription(description string, now time.Time) error {
	if description == "" {
		return nil
	}

	if now.IsZero() {
		return errors.New("now is not set")
	}

	b.description = description
	b.updatedAt = now

	return nil
}

func (b *Bag) ChangeStock(stock int32, now time.Time) error {
	if stock <= 0 {
		return errors.New("stock must be greater than 0")
	}

	if now.IsZero() {
		return errors.New("now is not set")
	}

	b.stock = stock
	b.updatedAt = now

	return nil
}

func (b *Bag) ChangePrice(originalPrice int64, discount int, now time.Time) error {
	if originalPrice <= 0 {
		return errors.New("original price must be greater than 0")
	}

	if discount < 30 || discount > 90 {
		return errors.New("discount must be between 30% and 90%")
	}

	if now.IsZero() {
		return errors.New("now is not set")
	}

	b.originalPrice = originalPrice
	b.finalPrice = calculateFinalPrice(b.originalPrice, discount)
	b.updatedAt = now

	return nil
}

func (b *Bag) ChangePickupWindow(pickupWindow PickupWindow, now time.Time) error {
	if now.IsZero() {
		return errors.New("now is not set")
	}

	b.pickupWindow = pickupWindow
	b.updatedAt = now

	return nil
}

func (b *Bag) ReplaceItems(items []Item, now time.Time) error {
	if len(items) <= 0 {
		return errors.New("at least one item is required")
	}

	if now.IsZero() {
		return errors.New("now is not set")
	}

	b.items = unique(items)
	b.updatedAt = now

	return nil
}

func (b Bag) ID() string {
	return b.id
}

func (b Bag) StoreID() string {
	return b.storeID
}

func (b Bag) Title() string {
	return b.title
}

func (b Bag) Description() string {
	return b.description
}

func (b Bag) Stock() int32 {
	return b.stock
}

func (b Bag) OriginalPrice() int64 {
	return b.originalPrice
}

func (b Bag) FinalPrice() int64 {
	return b.finalPrice
}

func (b Bag) Status() Status {
	return b.status
}

func (b Bag) PickupWindow() PickupWindow {
	return b.pickupWindow
}

func (b Bag) Items() []Item {
	items := make([]Item, len(b.items))
	copy(items, b.items)

	return items
}

func (b Bag) CreatedAt() time.Time {
	return b.createdAt
}

func (b Bag) UpdatedAt() time.Time {
	return b.updatedAt
}
