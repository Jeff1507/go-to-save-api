package repository

import (
	"database/sql"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/bag"
	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
)

type pickupWindowRow struct {
	Start time.Time
	End   time.Time
}

type itemRow struct {
	ID       string
	Name     string
	ImageKey sql.NullString
	Quantity int32
}

type bagRow struct {
	ID            string
	StoreID       string
	Title         string
	Description   sql.NullString
	Stock         int32
	OriginalPrice int64
	FinalPrice    int64
	Status        string
	PickupWindow  pickupWindowRow
	Items         []itemRow
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func toRow(b *bag.Bag) bagRow {
	var description sql.NullString
	if b.Description() != "" {
		description = sql.NullString{String: b.Description(), Valid: true}
	}

	pickupWindowRow := pickupWindowRow{
		Start: b.PickupWindow().Start(),
		End:   b.PickupWindow().End(),
	}

	itemsRow := make([]itemRow, len(b.Items()))
	for i, item := range b.Items() {
		var imageKey sql.NullString
		if item.ImageKey() != "" {
			imageKey = sql.NullString{String: item.ImageKey(), Valid: true}
		}

		itemsRow[i] = itemRow{
			ID:       item.ID(),
			Name:     item.Name(),
			ImageKey: imageKey,
			Quantity: item.Quantity(),
		}
	}

	return bagRow{
		ID:            b.ID(),
		StoreID:       b.StoreID(),
		Title:         b.Title(),
		Description:   description,
		Stock:         b.Stock(),
		OriginalPrice: b.OriginalPrice(),
		FinalPrice:    b.FinalPrice(),
		Status:        b.Status().String(),
		PickupWindow:  pickupWindowRow,
		Items:         itemsRow,
		CreatedAt:     b.CreatedAt(),
		UpdatedAt:     b.UpdatedAt(),
	}
}

func toEntity(r bagRow) *bag.Bag {
	status, _ := bag.ParseStatus(r.Status)

	pickupWindow, _ := pickup.NewWindow(r.PickupWindow.Start, r.PickupWindow.End)

	bagItems := make([]bag.Item, 0, len(r.Items))
	for _, item := range r.Items {
		bagItem, _ := bag.NewItem(item.ID, item.Name, item.ImageKey.String, item.Quantity)
		bagItems = append(bagItems, bagItem)
	}

	return bag.From(
		r.ID,
		r.StoreID,
		r.Title,
		r.Description.String,
		r.Stock,
		r.OriginalPrice,
		r.FinalPrice,
		status,
		pickupWindow,
		bagItems,
		r.CreatedAt,
		r.UpdatedAt,
	)
}
