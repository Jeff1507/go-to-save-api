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

type UpdateProfileAvatarInput struct {
	Content     io.Reader
	ContentType string
}

type UpdateProfileInput struct {
	ID        string
	FirstName string
	LastName  string
	Avatar    *UpdateProfileAvatarInput
}

type UpdateProfile struct {
	customerRepo customer.Repository
	storage      storage.Storage
	clock        clock.Clock
}

func NewUpdateProfile(
	customerRepo customer.Repository,
	storage storage.Storage,
	clock clock.Clock,
) *UpdateProfile {
	return &UpdateProfile{
		customerRepo: customerRepo,
		storage:      storage,
		clock:        clock,
	}
}

func (uc *UpdateProfile) Execute(ctx context.Context, input UpdateProfileInput) error {
	now := uc.clock.Now()

	cust, err := uc.customerRepo.Find(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("find customer by id: %w", err)
	}

	if err := cust.ChangeName(input.FirstName, input.LastName, now); err != nil {
		return fmt.Errorf("change first and last name: %w", err)
	}

	avatarChanged := input.Avatar != nil
	oldAvatarKey := cust.AvatarKey()

	var avatarKey string
	if avatarChanged {
		avatarKey, err = uc.storage.Upload(
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
	}

	if err := cust.ChangeAvatar(avatarKey, now); err != nil {
		uc.storage.Delete(avatarKey)
		return fmt.Errorf("change avatar: %w", err)
	}

	if err := uc.customerRepo.Update(ctx, cust); err != nil {
		uc.storage.Delete(avatarKey)
		return fmt.Errorf("update customer: %w", err)
	}

	if avatarChanged {
		uc.storage.Delete(oldAvatarKey)
	}

	return nil
}
