package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/customer"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
)

type RemoveAvatar struct {
	customerRepo customer.Repository
	storage      storage.Storage
	clock        clock.Clock
}

func NewRemoveAvatar(
	customerRepo customer.Repository,
	storage storage.Storage,
	clock clock.Clock,
) *RemoveAvatar {
	return &RemoveAvatar{
		customerRepo: customerRepo,
		storage:      storage,
		clock:        clock,
	}
}

func (uc *RemoveAvatar) Execute(ctx context.Context, id string) error {
	now := uc.clock.Now()

	cust, err := uc.customerRepo.Find(ctx, id)
	if err != nil {
		return fmt.Errorf("find customer by id: %w", err)
	}

	oldAvatarKey := cust.AvatarKey()
	if oldAvatarKey == "" {
		return nil
	}

	if err := cust.RemoveAvatar(now); err != nil {
		return fmt.Errorf("remove avatar: %w", err)
	}

	if err := uc.customerRepo.Update(ctx, cust); err != nil {
		return fmt.Errorf("update customer: %w", err)
	}

	uc.storage.Delete(oldAvatarKey)

	return nil
}
