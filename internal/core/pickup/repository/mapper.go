package repository

import (
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
)

type pickupTemplateRow struct {
	ID      string
	StoreID string
	Name    string
	Start   time.Time
	End     time.Time
}

func toEntity(r pickupTemplateRow) *pickup.Template {
	start, _ := pickup.NewTime(r.Start.Hour(), r.Start.Minute())
	end, _ := pickup.NewTime(r.End.Hour(), r.End.Minute())

	return pickup.TemplateFrom(r.ID, r.StoreID, r.Name, start, end)
}
