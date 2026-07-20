package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag/repository"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type closeBagDeps struct {
	clk  clock.Clock
	repo bag.Repository
	uc   *usecase.CloseBag
}

func newCloseBagDeps() *closeBagDeps {
	clk := clock.NewFakeClock()
	repo := repository.NewInMemRepo()
	uc := usecase.NewCloseBag(clk, repo)

	return &closeBagDeps{
		clk:  clk,
		repo: repo,
		uc:   uc,
	}
}

func TestCloseBag_Success(t *testing.T) {
	deps := newCloseBagDeps()

	now := deps.clk.Now()

	pickupStart := time.Date(2025, time.June, 1, 14, 0, 0, 0, time.UTC)
	pickupEnd := pickupStart.Add(time.Hour)

	win, err := bag.NewPickupWindow(pickupStart, pickupEnd)
	require.NoError(t, err)

	item1, err := bag.NewItem("prod-id", "prod-name", "", 1)
	require.NoError(t, err)

	b, err := bag.New("id", "store-id", "XDD", "", 2, 100, 30, win, []bag.Item{item1}, pickupStart)
	require.NoError(t, err)

	err = deps.repo.Save(context.Background(), b)
	require.NoError(t, err)

	err = deps.uc.Execute(context.Background(), b.ID())
	require.NoError(t, err)

	closed, err := deps.repo.Find(context.Background(), b.ID())
	require.NoError(t, err)

	assert.Equal(t, bag.Closed, closed.Status())
	assert.Equal(t, now, closed.UpdatedAt())
}

func TestCloseBag_NotFoundError(t *testing.T) {
	deps := newCloseBagDeps()

	err := deps.uc.Execute(context.Background(), "not-found-id")
	require.Error(t, err)

	assert.ErrorIs(t, err, bag.ErrNotFound)
}

func TestCloseBag_InvalidOperationError(t *testing.T) {
	deps := newCloseBagDeps()

	date := time.Date(2025, time.June, 1, 14, 0, 0, 0, time.UTC)

	b := bag.From(
		"id",
		"store-id",
		"Sacola",
		"Descrição",
		1,
		500,
		250,
		bag.Closed,
		bag.PickupWindow{},
		[]bag.Item{},
		date,
		date,
	)
	err := deps.repo.Save(context.Background(), b)
	require.NoError(t, err)

	err = deps.uc.Execute(context.Background(), b.ID())
	require.Error(t, err)

	assert.ErrorIs(t, err, bag.ErrInvalidOperation)
}
