package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/order"
	"github.com/Jeff1507/go-to-save-api/internal/core/review"
	"github.com/Jeff1507/go-to-save-api/internal/core/store"
	"github.com/Jeff1507/go-to-save-api/internal/infra/transaction"
	"github.com/Jeff1507/go-to-save-api/internal/pkg/uuid"
)

type AddReviewInput struct {
	OrderID    string
	CustomerID string
	Rating     int16
	Comment    string
}

type AddReviewOutput struct {
	ID string
}

type AddReview struct {
	transactions transaction.Manager
	reviewRepo   review.Repository
	orderRepo    order.Repository
	storeRepo    store.Repository
	clock        clock.Clock
}

func NewAddReview(
	transactions transaction.Manager,
	reviewRepo review.Repository,
	orderRepo order.Repository,
	storeRepo store.Repository,
	clock clock.Clock,
) *AddReview {
	return &AddReview{
		transactions: transactions,
		reviewRepo:   reviewRepo,
		orderRepo:    orderRepo,
		storeRepo:    storeRepo,
		clock:        clock,
	}
}

func (uc *AddReview) Execute(ctx context.Context, input AddReviewInput) (AddReviewOutput, error) {
	now := uc.clock.Now()
	var output AddReviewOutput

	err := uc.transactions.Do(ctx, func(ctx context.Context) error {
		ord, err := uc.orderRepo.Find(ctx, input.OrderID)
		if err != nil {
			if errors.Is(err, order.ErrNotFound) {
				return review.ErrCannotReview
			}

			return fmt.Errorf("find order by id: %w", err)
		}

		if ord.CustomerID() != input.CustomerID {
			return review.ErrCannotReview
		}

		if !uc.canReview(ord, now) {
			return review.ErrCannotReview
		}

		rev, err := review.New(uuid.Generate(), ord.ID(), input.Rating, input.Comment, now)
		if err != nil {
			return fmt.Errorf("new review: %w", err)
		}

		if err := uc.reviewRepo.Save(ctx, rev); err != nil {
			return fmt.Errorf("save review: %w", err)
		}

		if err := uc.storeRepo.FindAndUpdate(ctx, ord.StoreID(), func(s *store.Store) error {
			if err := s.AddRating(rev.Rating()); err != nil {
				return fmt.Errorf("add rating: %w", err)
			}

			return nil
		}); err != nil {
			return fmt.Errorf("find and update store: %w", err)
		}

		output = AddReviewOutput{
			ID: rev.ID(),
		}

		return nil
	})

	return output, err
}

func (uc *AddReview) canReview(o *order.Order, now time.Time) bool {
	const reviewWindow = time.Hour * 7 * 24

	if o.Status() != order.Completed {
		return false
	}

	deadline := o.CompletedAt().Add(reviewWindow)
	if o.CompletedAt().IsZero() || !now.Before(deadline) {
		return false
	}

	return true
}
