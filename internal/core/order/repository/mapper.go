package repository

import (
	"database/sql"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/order"
	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
)

type itemRow struct {
	BagID     string
	UnitPrice int64
	Quantity  int32
}

type orderRow struct {
	ID               string
	CustomerID       string
	StoreID          string
	TotalPrice       int64
	ConfirmationCode string
	Status           string
	PickupStart      time.Time
	PickupEnd        time.Time
	Items            []itemRow
	CompletedAt      sql.NullTime
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func toRow(o *order.Order) orderRow {
	var completedAt sql.NullTime
	if !o.CompletedAt().IsZero() {
		completedAt = sql.NullTime{Time: o.CompletedAt(), Valid: true}
	}

	itemsRow := make([]itemRow, len(o.Items()))
	for i, item := range o.Items() {
		itemsRow[i] = itemRow{
			BagID:     item.BagID(),
			UnitPrice: item.UnitPrice(),
			Quantity:  item.Quantity(),
		}
	}

	return orderRow{
		ID:               o.ID(),
		CustomerID:       o.CustomerID(),
		StoreID:          o.StoreID(),
		TotalPrice:       o.TotalPrice(),
		ConfirmationCode: o.ConfirmationCode(),
		Status:           o.Status().String(),
		PickupStart:      o.PickupWindow().Start(),
		PickupEnd:        o.PickupWindow().End(),
		Items:            itemsRow,
		CompletedAt:      completedAt,
		CreatedAt:        o.CreatedAt(),
		UpdatedAt:        o.UpdatedAt(),
	}
}

func toEntity(r orderRow) *order.Order {
	status, _ := order.ParseStatus(r.Status)

	pickupWindow, _ := pickup.NewWindow(r.PickupStart, r.PickupEnd)

	orderItems := make([]order.Item, len(r.Items))
	for i, item := range r.Items {
		orderItem, _ := order.NewItem(item.BagID, item.UnitPrice, item.Quantity)
		orderItems[i] = orderItem
	}

	return order.From(
		r.ID,
		r.CustomerID,
		r.StoreID,
		r.TotalPrice,
		r.ConfirmationCode,
		status,
		pickupWindow,
		orderItems,
		r.CompletedAt.Time,
		r.CreatedAt,
		r.UpdatedAt,
	)
}
