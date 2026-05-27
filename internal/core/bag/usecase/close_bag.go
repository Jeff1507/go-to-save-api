package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag"
)

type CloseBag struct {
	clock   clock.Clock
	bagRepo bag.Repository
}

func NewCloseBag(clock clock.Clock, bagRepo bag.Repository) *CloseBag {
	return &CloseBag{
		clock:   clock,
		bagRepo: bagRepo,
	}
}

func (uc *CloseBag) Execute(ctx context.Context, id string) error {
	now := uc.clock.Now()

	b, err := uc.bagRepo.Find(ctx, id)
	if err != nil {
		if errors.Is(err, bag.ErrNotFound) {
			return fmt.Errorf("find bag by id: %w", err)
		}

		return fmt.Errorf("find bag by id: %v", err)
	}

	if err := b.Close(now); err != nil {
		if errors.Is(err, bag.ErrInvalidOperation) {
			return fmt.Errorf("close bag: %w", err)
		}

		return fmt.Errorf("close bag: %v", err)
	}

	if err := uc.bagRepo.Update(ctx, b); err != nil {
		return fmt.Errorf("update bag: %v", err)
	}

	return nil
}
