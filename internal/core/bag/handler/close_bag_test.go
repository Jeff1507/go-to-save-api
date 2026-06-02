package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag/handler"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag/repository"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/validation"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newCloseBagReq(id string) *http.Request {
	escapedID := url.PathEscape(id)
	req := httptest.NewRequest(http.MethodPatch, "/bags/"+escapedID+"/close", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)

	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestCloseBag_Success(t *testing.T) {
	clk := clock.NewFakeClock()
	repo := repository.NewInMemRepo()
	uc := usecase.NewCloseBag(clk, repo)
	v := validation.NewValidator()
	h := handler.NewCloseBag(v, uc)

	now := clk.Now()

	win, err := bag.NewPickupWindow(now, now.Add(time.Hour))
	require.NoError(t, err)

	item1, err := bag.NewItem("prod-id", 3)
	require.NoError(t, err)

	b, err := bag.New(
		"a3f8c2d1-1234-4abc-89ef-123456789abc",
		"store-id",
		"Sacola mista",
		"...",
		1, 100, 30,
		win,
		[]bag.Item{item1},
		now,
	)
	require.NoError(t, err)

	err = repo.Save(context.Background(), b)
	require.NoError(t, err)

	req := newCloseBagReq(b.ID())
	rr := httptest.NewRecorder()

	h.Handle(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestCloseBag_ValidationError(t *testing.T) {
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
			clk := clock.NewFakeClock()
			repo := repository.NewInMemRepo()
			uc := usecase.NewCloseBag(clk, repo)
			v := validation.NewValidator()
			h := handler.NewCloseBag(v, uc)

			req := newCloseBagReq(tt.id)
			rr := httptest.NewRecorder()

			h.Handle(rr, req)

			assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		})
	}
}
