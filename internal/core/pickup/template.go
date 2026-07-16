package pickup

import (
	"errors"
	"time"
)

type Template struct {
	id      string
	storeID string
	name    string
	start   Time
	end     Time
}

func validateInterval(start, end time.Duration) error {
	const minInterval = time.Minute * 30 // 30 min
	const maxInterval = time.Hour * 6    // 06 hrs

	interval := end - start

	if end < start { // Virada de dia. Ex: 22:00 - 03:00
		interval = (end + time.Hour*24) - start
	}

	if interval < minInterval || interval > maxInterval {
		return ErrInvalidInterval
	}

	return nil
}

func NewTemplate(id, storeID, name string, start, end Time) (*Template, error) {
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

	return &Template{
		id:      id,
		storeID: storeID,
		name:    name,
		start:   start,
		end:     end,
	}, nil
}

func TemplateFrom(id, storeID, name string, start, end Time) *Template {
	return &Template{
		id:      id,
		storeID: storeID,
		name:    name,
		start:   start,
		end:     end,
	}
}

func (t *Template) ChangeName(name string) error {
	if name == "" {
		return errors.New("name is required")
	}

	t.name = name

	return nil
}

func (t *Template) ChangeTime(start, end Time) error {
	if err := validateInterval(start.Duration(), end.Duration()); err != nil {
		return err
	}

	t.start = start
	t.end = end

	return nil
}

func (t Template) ID() string {
	return t.id
}

func (t Template) StoreID() string {
	return t.storeID
}

func (t Template) Name() string {
	return t.name
}

func (t Template) Start() Time {
	return t.start
}

func (t Template) End() Time {
	return t.end
}
