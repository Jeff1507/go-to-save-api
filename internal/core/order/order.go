package order

import (
	"errors"
	"fmt"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
)

type Status uint8

const (
	Unknown Status = iota
	Pending
	Completed
	Cancelled
	NoShow
)

func (s Status) String() string {
	switch s {
	case Pending:
		return "PENDING"

	case Completed:
		return "COMPLETED"

	case Cancelled:
		return "CANCELLED"

	case NoShow:
		return "NO_SHOW"

	default:
		return "UNKNOWN"
	}
}

func ParseStatus(str string) (Status, error) {
	switch str {
	case "PENDING":
		return Pending, nil

	case "COMPLETED":
		return Completed, nil

	case "CANCELLED":
		return Cancelled, nil

	case "NO_SHOW":
		return NoShow, nil

	default:
		return Unknown, fmt.Errorf("unknown status: %q", str)
	}
}

type Order struct {
	id               string
	customerID       string
	storeID          string
	totalPrice       int64 // Preço em centavos
	confirmationCode string
	status           Status
	pickupWindow     pickup.Window
	items            []Item
	completedAt      time.Time
	createdAt        time.Time
	updatedAt        time.Time
}

func unique(items []Item) ([]Item, error) {
	seen := make(map[string]bool, len(items))
	for _, item := range items {
		if seen[item.bagID] {
			return nil, ErrDuplicateItems
		}

		seen[item.bagID] = true
	}

	return items, nil
}

func New(
	id string,
	customerID string,
	storeID string,
	pickupWindow pickup.Window,
	items []Item,
	now time.Time,
) (*Order, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	if customerID == "" {
		return nil, errors.New("customer id is required")
	}

	if storeID == "" {
		return nil, errors.New("store id is required")
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

	confirmationCode, err := GenerateConfirmationCode()
	if err != nil {
		return nil, fmt.Errorf("generate confirmation code: %w", err)
	}

	uniqueItems, err := unique(items)
	if err != nil {
		return nil, err
	}

	var totalPrice int64
	for _, item := range uniqueItems {
		totalPrice += item.TotalPrice()
	}

	return &Order{
		id:               id,
		customerID:       customerID,
		storeID:          storeID,
		totalPrice:       totalPrice,
		confirmationCode: confirmationCode,
		status:           Pending,
		pickupWindow:     pickupWindow,
		items:            uniqueItems,
		createdAt:        now,
		updatedAt:        now,
	}, nil
}

func From(
	id string,
	customerID string,
	storeID string,
	totalPrice int64,
	confirmationCode string,
	status Status,
	pickupWindow pickup.Window,
	items []Item,
	completedAt time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) *Order {
	return &Order{
		id:               id,
		customerID:       customerID,
		storeID:          storeID,
		totalPrice:       totalPrice,
		confirmationCode: confirmationCode,
		status:           status,
		pickupWindow:     pickupWindow,
		items:            items,
		completedAt:      completedAt,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
	}
}

func (o Order) IsPending() bool {
	return o.status == Pending
}

func (o Order) IsExpired(now time.Time) bool {
	return o.PickupWindow().IsExpired(now)
}

func (o *Order) Complete(confirmationCode string, now time.Time) error {
	if !o.IsPending() {
		return ErrInvalidState
	}

	if o.confirmationCode != confirmationCode {
		return ErrInvalidCode
	}

	o.status = Completed
	o.completedAt = now
	o.updatedAt = now

	return nil
}

func (o *Order) Cancel(now time.Time) error {
	if !o.IsPending() {
		return ErrInvalidState
	}

	o.status = Cancelled
	o.updatedAt = now

	return nil
}

func (o *Order) NoShow(now time.Time) error {
	if !o.IsPending() {
		return ErrInvalidState
	}

	o.status = NoShow
	o.updatedAt = now

	return nil
}

func (o Order) ID() string {
	return o.id
}

func (o Order) CustomerID() string {
	return o.customerID
}

func (o Order) StoreID() string {
	return o.storeID
}

func (o Order) TotalPrice() int64 {
	return o.totalPrice
}

func (o Order) ConfirmationCode() string {
	return o.confirmationCode
}

func (o Order) Status() Status {
	return o.status
}

func (o Order) PickupWindow() pickup.Window {
	return o.pickupWindow
}

func (o Order) Items() []Item {
	items := make([]Item, len(o.items))
	copy(items, o.items)

	return items
}

func (o Order) CompletedAt() time.Time {
	return o.completedAt
}

func (o Order) CreatedAt() time.Time {
	return o.createdAt
}

func (o Order) UpdatedAt() time.Time {
	return o.updatedAt
}
