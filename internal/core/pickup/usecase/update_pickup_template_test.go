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

func newUpdatePickupTemplateDeps() (pickup.Repository, *usecase.UpdatePickupTemplate) {
	repo := repository.NewInMemRepo()
	uc := usecase.NewUpdatePickupTemplate(repo)

	return repo, uc
}

func createPickupTemplateForUpdate(t *testing.T, repo pickup.Repository) *pickup.Template {
	t.Helper()

	start, err := pickup.NewTime(18, 0)
	require.NoError(t, err)

	end, err := pickup.NewTime(22, 0)
	require.NoError(t, err)

	templ, err := pickup.NewTemplate("id", "store-id", "Fechamento", start, end)
	require.NoError(t, err)

	err = repo.Save(context.Background(), templ)
	require.NoError(t, err)

	return templ
}

func TestUpdatePickupTemplate_Success(t *testing.T) {
	repo, uc := newUpdatePickupTemplateDeps()

	templ := createPickupTemplateForUpdate(t, repo)

	input := usecase.UpdatePickupTemplateInput{
		ID:   templ.ID(),
		Name: "Almoço",
		Start: usecase.UpdatePickupTemplateTimeInput{
			Hour:   12,
			Minute: 0,
		},
		End: usecase.UpdatePickupTemplateTimeInput{
			Hour:   15,
			Minute: 30,
		},
	}

	err := uc.Execute(context.Background(), input)
	require.NoError(t, err)

	updated, err := repo.Find(context.Background(), templ.ID())
	require.NoError(t, err)

	assert.Equal(t, input.Name, updated.Name())
	assert.Equal(t, input.Start.Hour, updated.Start().Hour())
	assert.Equal(t, input.Start.Minute, updated.Start().Minute())
	assert.Equal(t, input.End.Hour, updated.End().Hour())
	assert.Equal(t, input.End.Minute, updated.End().Minute())
}

func TestUpdatePickupTemplate_SameName_Success(t *testing.T) {
	repo, uc := newUpdatePickupTemplateDeps()

	templ := createPickupTemplateForUpdate(t, repo)

	input := usecase.UpdatePickupTemplateInput{
		ID:   templ.ID(),
		Name: templ.Name(),
		Start: usecase.UpdatePickupTemplateTimeInput{
			Hour:   12,
			Minute: 0,
		},
		End: usecase.UpdatePickupTemplateTimeInput{
			Hour:   15,
			Minute: 30,
		},
	}

	err := uc.Execute(context.Background(), input)
	require.NoError(t, err)

	updated, err := repo.Find(context.Background(), templ.ID())
	require.NoError(t, err)

	assert.Equal(t, templ.Name(), updated.Name())
	assert.Equal(t, input.Start.Hour, updated.Start().Hour())
	assert.Equal(t, input.Start.Minute, updated.Start().Minute())
	assert.Equal(t, input.End.Hour, updated.End().Hour())
	assert.Equal(t, input.End.Minute, updated.End().Minute())
}

func TestUpdatePickupTemplate_NotFoundError(t *testing.T) {
	_, uc := newUpdatePickupTemplateDeps()

	input := usecase.UpdatePickupTemplateInput{
		ID:   "not-found-id",
		Name: "Not found",
		Start: usecase.UpdatePickupTemplateTimeInput{
			Hour:   20,
			Minute: 0,
		},
		End: usecase.UpdatePickupTemplateTimeInput{
			Hour:   22,
			Minute: 30,
		},
	}

	err := uc.Execute(context.Background(), input)
	require.Error(t, err)

	assert.ErrorIs(t, err, pickup.ErrTemplateNotFound)
}

func TestUpdatePickupTemplate_AlreadyExistsError(t *testing.T) {
	repo, uc := newUpdatePickupTemplateDeps()

	start, err := pickup.NewTime(15, 0)
	require.NoError(t, err)

	end, err := pickup.NewTime(18, 0)
	require.NoError(t, err)

	existed, err := pickup.NewTemplate("existed-id", "store-id", "Fim do dia", start, end)
	require.NoError(t, err)

	err = repo.Save(context.Background(), existed)
	require.NoError(t, err)

	templ := createPickupTemplateForUpdate(t, repo)

	input := usecase.UpdatePickupTemplateInput{
		ID:   templ.ID(),
		Name: existed.Name(),
		Start: usecase.UpdatePickupTemplateTimeInput{
			Hour:   20,
			Minute: 0,
		},
		End: usecase.UpdatePickupTemplateTimeInput{
			Hour:   22,
			Minute: 30,
		},
	}

	err = uc.Execute(context.Background(), input)
	require.Error(t, err)

	assert.ErrorIs(t, err, pickup.ErrTemplateAlreadyExists)

	notUpdated, err := repo.Find(context.Background(), templ.ID())
	require.NoError(t, err)

	assert.Equal(t, templ.Name(), notUpdated.Name())
	assert.Equal(t, templ.Start(), notUpdated.Start())
	assert.Equal(t, templ.End(), notUpdated.End())
}

func TestUpdatePickupTemplate_InvalidIntervalError(t *testing.T) {
	tests := []struct {
		name  string
		start usecase.UpdatePickupTemplateTimeInput
		end   usecase.UpdatePickupTemplateTimeInput
	}{
		{
			name:  "start equal to end",
			start: usecase.UpdatePickupTemplateTimeInput{22, 0},
			end:   usecase.UpdatePickupTemplateTimeInput{22, 0},
		},
		{
			name:  "interval is less than 30 min",
			start: usecase.UpdatePickupTemplateTimeInput{0, 0},
			end:   usecase.UpdatePickupTemplateTimeInput{0, 29},
		},
		{
			name:  "interval is greater than 6 hours",
			start: usecase.UpdatePickupTemplateTimeInput{0, 0},
			end:   usecase.UpdatePickupTemplateTimeInput{6, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, uc := newUpdatePickupTemplateDeps()

			templ := createPickupTemplateForUpdate(t, repo)

			input := usecase.UpdatePickupTemplateInput{
				ID:    templ.ID(),
				Name:  "Madrugada",
				Start: tt.start,
				End:   tt.end,
			}

			err := uc.Execute(context.Background(), input)
			require.Error(t, err)

			assert.ErrorIs(t, err, pickup.ErrInvalidInterval)

			notUpdated, err := repo.Find(context.Background(), templ.ID())
			require.NoError(t, err)

			assert.Equal(t, templ.Name(), notUpdated.Name())
			assert.Equal(t, templ.Start(), notUpdated.Start())
			assert.Equal(t, templ.End(), notUpdated.End())
		})
	}
}
