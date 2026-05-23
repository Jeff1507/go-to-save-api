package usecase_test

import (
	"context"
	"testing"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/repository"
	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/pkg/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newCreatePickupTemplateDeps() (pickup.Repository, *usecase.CreatePickupTemplate) {
	repo := repository.NewInMemRepo()
	uc := usecase.NewCreatePickupTemplate(repo)

	return repo, uc
}

func TestCreatePickupTemplate_Success(t *testing.T) {
	repo, uc := newCreatePickupTemplateDeps()

	input := usecase.CreatePickupTemplateInput{
		StoreID: "store-id",
		Name:    "Madrugada",
		Start: usecase.CreatePickupTemplateTimeInput{
			Hour:   22,
			Minute: 0,
		},
		End: usecase.CreatePickupTemplateTimeInput{
			Hour:   3,
			Minute: 30,
		},
	}

	output, err := uc.Execute(context.Background(), input)
	require.NoError(t, err)

	templ, err := repo.Find(context.Background(), output.ID)
	require.NoError(t, err)

	assert.True(t, uuid.IsValid(templ.ID()), "pickup template id should be an uuid")
	assert.Equal(t, input.StoreID, templ.StoreID())
	assert.Equal(t, input.Name, templ.Name())
	assert.Equal(t, input.Start.Hour, templ.Start().Hour())
	assert.Equal(t, input.Start.Minute, templ.Start().Minute())
	assert.Equal(t, input.End.Hour, templ.End().Hour())
	assert.Equal(t, input.End.Minute, templ.End().Minute())
}

func TestCreatePickupTemplate_AlreadyExistsError(t *testing.T) {
	repo, uc := newCreatePickupTemplateDeps()

	start, err := pickup.NewTime(18, 0)
	require.NoError(t, err)

	end, err := pickup.NewTime(20, 0)
	require.NoError(t, err)

	existed, err := pickup.NewTemplate("id", "store-id", "Fechamento", start, end)
	require.NoError(t, err)

	err = repo.Save(context.Background(), existed)
	require.NoError(t, err)

	input := usecase.CreatePickupTemplateInput{
		StoreID: existed.StoreID(),
		Name:    existed.Name(),
		Start: usecase.CreatePickupTemplateTimeInput{
			Hour:   22,
			Minute: 0,
		},
		End: usecase.CreatePickupTemplateTimeInput{
			Hour:   23,
			Minute: 59,
		},
	}

	output, err := uc.Execute(context.Background(), input)
	require.Error(t, err)

	assert.ErrorIs(t, err, pickup.ErrTemplateAlreadyExists)

	_, err = repo.Find(context.Background(), output.ID)
	require.Error(t, err)

	assert.ErrorIs(t, err, pickup.ErrTemplateNotFound)
}

func TestCreatePickupTemplate_InvalidIntervalError(t *testing.T) {
	tests := []struct {
		name  string
		start usecase.CreatePickupTemplateTimeInput
		end   usecase.CreatePickupTemplateTimeInput
	}{
		{
			name:  "start equal to end",
			start: usecase.CreatePickupTemplateTimeInput{22, 0},
			end:   usecase.CreatePickupTemplateTimeInput{22, 0},
		},
		{
			name:  "interval is less than 30 minutes",
			start: usecase.CreatePickupTemplateTimeInput{22, 0},
			end:   usecase.CreatePickupTemplateTimeInput{22, 29},
		},
		{
			name:  "interval is greater than 6 hours",
			start: usecase.CreatePickupTemplateTimeInput{22, 0},
			end:   usecase.CreatePickupTemplateTimeInput{4, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, uc := newCreatePickupTemplateDeps()

			input := usecase.CreatePickupTemplateInput{
				StoreID: "store-id",
				Name:    "Domingo a noite",
				Start:   tt.start,
				End:     tt.end,
			}

			output, err := uc.Execute(context.Background(), input)
			require.Error(t, err)

			assert.ErrorIs(t, err, pickup.ErrInvalidInterval)

			_, err = repo.Find(context.Background(), output.ID)
			require.Error(t, err)

			assert.ErrorIs(t, err, pickup.ErrTemplateNotFound)
		})
	}
}
