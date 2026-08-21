package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/customer"
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
	customerRepo customer.Repository
	storage      storage.Storage
	clock        clock.Clock
}

func NewUpdateAvatar(
	customerRepo customer.Repository,
	storage storage.Storage,
	clock clock.Clock,
) *UpdateAvatar {
	return &UpdateAvatar{
		customerRepo: customerRepo,
		storage:      storage,
		clock:        clock,
	}
}

func (uc *UpdateAvatar) Execute(ctx context.Context, input UpdateAvatarInput) error {
	now := uc.clock.Now()

	cust, err := uc.customerRepo.Find(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("find customer by id: %w", err)
	}

	oldAvatarKey := cust.AvatarKey()

	newAvatarKey, err := uc.storage.Upload(
		input.Avatar.Content,
		input.Avatar.ContentType,
		storage.AvatarsPath,
	)
	if err != nil {
		if errors.Is(err, storage.ErrContentTypeNotAllowed) {
			return customer.ErrInvalidImageFormat
		}

		return fmt.Errorf("upload avatar: %w", err)
	}

	if err := cust.ChangeAvatar(newAvatarKey, now); err != nil {
		uc.storage.Delete(newAvatarKey)
		return fmt.Errorf("change avatar: %w", err)
	}

	if err := uc.customerRepo.Update(ctx, cust); err != nil {
		uc.storage.Delete(newAvatarKey)
		return fmt.Errorf("update customer: %w", err)
	}

	uc.storage.Delete(oldAvatarKey)

	return nil
}
