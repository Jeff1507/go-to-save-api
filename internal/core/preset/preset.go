package preset

import (
	"errors"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
)

type Preset struct {
	id      string
	storeID string
	name    string
	start   pickup.Time
	end     pickup.Time
}

func validateInterval(start, end time.Duration) error {
	const minInterval = time.Minute * 30 // 30 min
	const maxInterval = time.Hour * 6    // 06 hrs

	interval := end - start

	if end <= start { // Virada de dia. Ex: 22:00 - 03:00
		interval = (end + time.Hour*24) - start
	}

	if interval < minInterval || interval > maxInterval {
		return ErrInvalidInterval
	}

	return nil
}

func New(id, storeID, name string, start, end pickup.Time) (*Preset, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	if storeID == "" {
		return nil, errors.New("store id is required")
	}

	if name == "" {
		return nil, errors.New("name is required")
	}

	if err := validateInterval(start.Duration(), end.Duration()); err != nil {
		return nil, err
	}

	return &Preset{
		id:      id,
		storeID: storeID,
		name:    name,
		start:   start,
		end:     end,
	}, nil
}

func From(id, storeID, name string, start, end pickup.Time) *Preset {
	return &Preset{
		id:      id,
		storeID: storeID,
		name:    name,
		start:   start,
		end:     end,
	}
}

func (p *Preset) ChangeName(name string) error {
	if name == "" {
		return errors.New("name is required")
	}

	p.name = name

	return nil
}

func (p *Preset) ChangePickupTime(start, end pickup.Time) error {
	if err := validateInterval(start.Duration(), end.Duration()); err != nil {
		return err
	}

	p.start = start
	p.end = end

	return nil
}

func (p Preset) ID() string {
	return p.id
}

func (p Preset) StoreID() string {
	return p.storeID
}

func (p Preset) Name() string {
	return p.name
}

func (p Preset) Start() pickup.Time {
	return p.start
}

func (p Preset) End() pickup.Time {
	return p.end
}
