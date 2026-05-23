package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/product"
	"github.com/Jeff1507/go-to-save-api/internal/core/product/handler"
	"github.com/Jeff1507/go-to-save-api/internal/core/product/repository"
	"github.com/Jeff1507/go-to-save-api/internal/core/product/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newDeleteProductReq(id string) *http.Request {
	escapedID := url.PathEscape(id)
	req := httptest.NewRequest(http.MethodDelete, "/products/"+escapedID, nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)

	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestDeleteProduct_Success(t *testing.T) {
	st := storage.NewFakeStorage()
	repo := repository.NewInMemRepo()
	uc := usecase.NewDeleteProduct(st, repo)
	h := handler.NewDeleteProduct(uc)
	date := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)

	prod, err := product.New(
		"a3f8c2d1-1234-4abc-89ef-123456789abc",
		"store-id",
		"Coxinha",
		"coxinha.png",
		date,
	)
	require.NoError(t, err)

	err = repo.Save(context.Background(), prod)
	require.NoError(t, err)

	req := newDeleteProductReq(prod.ID())
	rr := httptest.NewRecorder()

	h.Handle(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteProduct_ValidationError(t *testing.T) {
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
			st := storage.NewFakeStorage()
			repo := repository.NewInMemRepo()
			uc := usecase.NewDeleteProduct(st, repo)
			h := handler.NewDeleteProduct(uc)

			req := newDeleteProductReq(tt.id)
			rr := httptest.NewRecorder()

			h.Handle(rr, req)

			assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		})
	}
}
