package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag"
)

type CloseBagInput struct {
	ID      string
	StoreID string
}

type CloseBag struct {
	bagRepo bag.Repository
	clock   clock.Clock
}

func NewCloseBag(bagRepo bag.Repository, clock clock.Clock) *CloseBag {
	return &CloseBag{
		bagRepo: bagRepo,
		clock:   clock,
	}
}

func (uc *CloseBag) Execute(ctx context.Context, input CloseBagInput) error {
	now := uc.clock.Now()

	b, err := uc.bagRepo.Find(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("find bag by id: %w", err)
	}

	if b.StoreID() != input.StoreID {
		return bag.ErrNotFound
	}

	if err := b.Close(now); err != nil {
		return fmt.Errorf("close bag: %w", err)
	}

	if err := uc.bagRepo.Update(ctx, b); err != nil {
		return fmt.Errorf("update bag: %w", err)
	}

	return nil
}
