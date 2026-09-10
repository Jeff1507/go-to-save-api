package bag

import (
	"errors"
	"fmt"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
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
	pickupWindow  pickup.Window
	items         []Item
	createdAt     time.Time
	updatedAt     time.Time
}

func calculateFinalPrice(originalPrice int64, discount int) int64 {
	return originalPrice - (originalPrice * int64(discount) / 100)
}

func unique(items []Item) ([]Item, error) {
	seen := make(map[string]bool, len(items))
	for _, item := range items {
		if seen[item.name] {
			return nil, ErrDuplicateItems
		}

		seen[item.name] = true
	}

	return items, nil
}

func New(
	id string,
	storeID string,
	title string,
	description string,
	stock int32,
	originalPrice int64,
	discount int,
	pickupWindow pickup.Window,
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

	if pickupWindow.IsZero() {
		return nil, errors.New("pickup window is required")
	}

	if len(items) == 0 {
		return nil, errors.New("at least one item is required")
	}

	if now.IsZero() {
		return nil, errors.New("now is not set")
	}

	uniqueItems, err := unique(items)
	if err != nil {
		return nil, err
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
		items:         uniqueItems,
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
	pickupWindow pickup.Window,
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
	return b.PickupWindow().IsExpired(now)
}

func (b Bag) IsClosed() bool {
	return b.status == Closed
}

func (b *Bag) Close(now time.Time) error {
	if now.IsZero() {
		return errors.New("now is not set")
	}

	if b.IsClosed() {
		return ErrClosed
	}

	b.status = Closed
	b.updatedAt = now

	return nil
}

func (b *Bag) Reserve(quantity int32, now time.Time) error {
	if b.IsExpired(now) {
		return ErrExpired
	}

	if b.IsClosed() {
		return ErrClosed
	}

	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	if quantity > b.stock {
		return ErrInsufficientStock
	}

	if now.IsZero() {
		return errors.New("now is not set")
	}

	b.stock -= quantity
	b.updatedAt = now

	return nil
}

// Restock devolve quantidade ao estoque da sacola, tipicamente
// após cancelamento ou no-show de um pedido. Não valida status/expiração
// da sacola, pois a devolução deve ser possível mesmo após o fechamento.
func (b *Bag) Restock(quantity int32, now time.Time) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	if now.IsZero() {
		return errors.New("now is not set")
	}

	b.stock += quantity
	b.updatedAt = now

	return nil
}

// AdjustStock aplica um ajuste manual de estoque feito pela loja (ex: reposição,
// correção de contagem, remoção por perda). Delta pode ser positivo ou negativo.
// Diferente de Reserve (cliente) e Restock (sistema, após cancelamento), esta é
// uma ação administrativa e só é permitida em sacolas ainda ativas.
func (b *Bag) AdjustStock(delta int32, now time.Time) error {
	if now.IsZero() {
		return errors.New("now is not set")
	}

	if b.IsClosed() {
		return ErrClosed
	}

	newStock := b.stock + delta
	if newStock < 0 {
		return ErrInsufficientStock
	}

	b.stock = newStock
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
	if now.IsZero() {
		return errors.New("now is not set")
	}

	b.description = description
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

func (b *Bag) ChangePickupWindow(pickupWindow pickup.Window, now time.Time) error {
	if pickupWindow.IsZero() {
		return errors.New("pickup window is required")
	}

	if now.IsZero() {
		return errors.New("now is not set")
	}

	b.pickupWindow = pickupWindow
	b.updatedAt = now

	return nil
}

func (b *Bag) ReplaceItems(items []Item, now time.Time) error {
	if len(items) == 0 {
		return errors.New("at least one item is required")
	}

	if now.IsZero() {
		return errors.New("now is not set")
	}

	uniqueItems, err := unique(items)
	if err != nil {
		return err
	}

	b.items = uniqueItems
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

func (b Bag) PickupWindow() pickup.Window {
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
