package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/store"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
)

type AvatarInput struct {
	Content     io.Reader
	ContentType string
}

type UpdateAvatarInput struct {
	ID     string
	Avatar AvatarInput
}

type UpdateAvatar struct {
	storeRepo store.Repository
	storage   storage.Storage
	clock     clock.Clock
}

func NewUpdateAvatar(
	storeRepo store.Repository,
	storage storage.Storage,
	clock clock.Clock,
) *UpdateAvatar {
	return &UpdateAvatar{
		storeRepo: storeRepo,
		storage:   storage,
		clock:     clock,
	}
}

func (uc *UpdateAvatar) Execute(ctx context.Context, input UpdateAvatarInput) error {
	now := uc.clock.Now()

	s, err := uc.storeRepo.Find(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("find store by id: %w", err)
	}

	oldAvatarKey := s.AvatarKey()

	newAvatarKey, err := uc.storage.Upload(
		input.Avatar.Content,
		input.Avatar.ContentType,
		storage.AvatarsPath,
	)
	if err != nil {
		if errors.Is(err, storage.ErrContentTypeNotAllowed) {
			return store.ErrInvalidImageFormat
		}

		return fmt.Errorf("upload avatar: %w", err)
	}

	if err := s.ChangeAvatar(newAvatarKey, now); err != nil {
		uc.storage.Delete(newAvatarKey)
		return fmt.Errorf("change avatar: %w", err)
	}

	if err := uc.storeRepo.Update(ctx, s); err != nil {
		uc.storage.Delete(newAvatarKey)
		return fmt.Errorf("update store: %w", err)
	}

	uc.storage.Delete(oldAvatarKey)

	return nil
}
