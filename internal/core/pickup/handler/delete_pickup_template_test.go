package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/handler"
	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/repository"
	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newDeletePickupTemplateReq(id string) *http.Request {
	escapedID := url.PathEscape(id)
	req := httptest.NewRequest(http.MethodDelete, "/pickup-templates/"+escapedID, nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)

	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestDeletePickupTemplate_Success(t *testing.T) {
	repo := repository.NewInMemRepo()
	uc := usecase.NewDeletePickupTemplate(repo)
	h := handler.NewDeletePickupTemplate(uc)

	start, err := pickup.NewTime(20, 30)
	require.NoError(t, err)

	end, err := pickup.NewTime(22, 0)
	require.NoError(t, err)

	templ, err := pickup.NewTemplate(
		"a3f8c2d1-1234-4abc-89ef-123456789abc",
		"store-id",
		"Domingo a noite",
		start,
		end,
	)
	require.NoError(t, err)

	err = repo.Save(context.Background(), templ)
	require.NoError(t, err)

	req := newDeletePickupTemplateReq(templ.ID())
	rr := httptest.NewRecorder()

	h.Handle(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeletePickupTemplate_ValidationError(t *testing.T) {
	tests := []struct {
		name string
		id   string
	}{
		{
			name: "empty id",
		},
		{
			name: "id with white space",
			id:   "   ",
		},
		{
			name: "invalid id",
			id:   "invalid-id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repository.NewInMemRepo()
			uc := usecase.NewDeletePickupTemplate(repo)
			h := handler.NewDeletePickupTemplate(uc)

			req := newDeletePickupTemplateReq(tt.id)
			rr := httptest.NewRecorder()

			h.Handle(rr, req)

			assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		})
	}
}
