package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/store"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
)

type RemoveAvatar struct {
	storeRepo store.Repository
	storage   storage.Storage
	clock     clock.Clock
}

func NewRemoveAvatar(
	storeRepo store.Repository,
	storage storage.Storage,
	clock clock.Clock,
) *RemoveAvatar {
	return &RemoveAvatar{
		storeRepo: storeRepo,
		storage:   storage,
		clock:     clock,
	}
}

func (uc *RemoveAvatar) Execute(ctx context.Context, id string) error {
	now := uc.clock.Now()

	s, err := uc.storeRepo.Find(ctx, id)
	if err != nil {
		return fmt.Errorf("find store by id: %w", err)
	}

	oldAvatarKey := s.AvatarKey()
	if oldAvatarKey == "" {
		return nil
	}

	if err := s.RemoveAvatar(now); err != nil {
		return fmt.Errorf("remove avatar: %w", err)
	}

	if err := uc.storeRepo.Update(ctx, s); err != nil {
		return fmt.Errorf("update store: %w", err)
	}

	uc.storage.Delete(oldAvatarKey)

	return nil
}
