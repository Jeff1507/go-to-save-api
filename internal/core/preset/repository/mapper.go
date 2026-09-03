package repository

import (
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
	"github.com/Jeff1507/go-to-save-api/internal/core/preset"
)

type presetRow struct {
	ID      string
	StoreID string
	Name    string
	Start   time.Time
	End     time.Time
}

func toEntity(r presetRow) *preset.Preset {
	start, _ := pickup.NewTime(r.Start.Hour(), r.Start.Minute())
	end, _ := pickup.NewTime(r.End.Hour(), r.End.Minute())

	return preset.From(r.ID, r.StoreID, r.Name, start, end)
}
