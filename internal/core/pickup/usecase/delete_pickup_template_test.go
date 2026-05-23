package usecase_test

import (
	"context"
	"testing"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/repository"
	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newDeletePickupTemplateDeps() (pickup.Repository, *usecase.DeletePickupTemplate) {
	repo := repository.NewInMemRepo()
	uc := usecase.NewDeletePickupTemplate(repo)

	return repo, uc
}

func TestDeletePickupTemplate_Success(t *testing.T) {
	repo, uc := newDeletePickupTemplateDeps()

	start, err := pickup.NewTime(12, 0)
	require.NoError(t, err)

	end, err := pickup.NewTime(13, 0)
	require.NoError(t, err)

	templ, err := pickup.NewTemplate("id", "store-id", "Depois do Almoço", start, end)
	require.NoError(t, err)

	err = repo.Save(context.Background(), templ)
	require.NoError(t, err)

	err = uc.Execute(context.Background(), templ.ID())
	require.NoError(t, err)

	_, err = repo.Find(context.Background(), templ.ID())
	require.Error(t, err)

	assert.ErrorIs(t, err, pickup.ErrTemplateNotFound)
}

func TestDeletePickupTemplate_NotFound_Success(t *testing.T) {
	_, uc := newDeletePickupTemplateDeps()

	err := uc.Execute(context.Background(), "not-found-id")
	require.NoError(t, err)
}
