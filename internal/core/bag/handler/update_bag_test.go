package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag/handler"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag/repository"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/core/product"
	productRepository "github.com/Jeff1507/go-to-save-api/internal/core/product/repository"
	"github.com/Jeff1507/go-to-save-api/internal/infra/validation"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newUpdateBagReq(id string, body io.Reader) *http.Request {
	escapedID := url.PathEscape(id)
	req := httptest.NewRequest(http.MethodPut, "/bags/"+escapedID, body)
	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)

	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func createBagToUpdate(
	t *testing.T,
	productRepo product.Repository,
	bagRepo bag.Repository,
	pickupStart time.Time,
	pickupEnd time.Time,
) *bag.Bag {
	t.Helper()

	date := time.Date(2024, time.April, 1, 0, 0, 0, 0, time.UTC)

	prod, err := product.New(
		"550e8400-e29b-41d4-a716-446655440000",
		"store-id",
		"Coxinha",
		"coxinha.png",
		date,
	)
	require.NoError(t, err)

	err = productRepo.Save(context.Background(), prod)
	require.NoError(t, err)

	win, err := bag.NewPickupWindow(pickupStart, pickupEnd)
	require.NoError(t, err)

	item, err := bag.NewItem(prod.ID(), 2)
	require.NoError(t, err)

	b, err := bag.New(
		"a3f8c2d1-1234-4abc-89ef-123456789abc",
		"store-id",
		"XDD",
		"", 1, 500, 50,
		win,
		[]bag.Item{item},
		pickupStart,
	)
	require.NoError(t, err)

	err = bagRepo.Save(context.Background(), b)
	require.NoError(t, err)

	return b
}

func TestUpdateBag_Success(t *testing.T) {
	clk := clock.NewFakeClock()
	bagRepo := repository.NewInMemRepo()
	productRepo := productRepository.NewInMemRepo()
	uc := usecase.NewUpdateBag(clk, bagRepo, productRepo)
	v := validation.NewValidator()
	h := handler.NewUpdateBag(v, uc)

	now := clk.Now()

	pickupStart := now.Add(time.Hour * 24)
	pickupEnd := pickupStart.Add(time.Hour)

	b := createBagToUpdate(t, productRepo, bagRepo, pickupStart, pickupEnd)

	body := handler.UpdateBagRequestBody{
		Title:         "Sacola editada",
		Description:   "descrição editada...",
		Stock:         3,
		OriginalPrice: 1000,
		Discount:      40,
		PickupWindow: handler.UpdateBagPickupWindowRequestBody{
			Start: pickupStart,
			End:   pickupEnd,
		},
		Items: []handler.UpdateBagItemRequestBody{
			{ID: b.Items()[0].ID(), Quantity: 2},
		},
	}

	bt, err := json.Marshal(body)
	require.NoError(t, err)

	req := newUpdateBagReq(b.ID(), bytes.NewReader(bt))
	rr := httptest.NewRecorder()

	h.Handle(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestUpdateBag_ValidationError(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	validID := "a3f8c2d1-1234-4abc-89ef-123456789abc"

	validBody := handler.UpdateBagRequestBody{
		Title:         "Sacola editada",
		Description:   "descrição editada...",
		Stock:         3,
		OriginalPrice: 1000,
		Discount:      40,
		PickupWindow: handler.UpdateBagPickupWindowRequestBody{
			Start: now,
			End:   now.Add(time.Hour),
		},
		Items: []handler.UpdateBagItemRequestBody{
			{ID: "550e8400-e29b-41d4-a716-446655440000", Quantity: 2},
		},
	}

	tests := []struct {
		name    string
		id      string
		reqBody handler.UpdateBagRequestBody
	}{
		{
			name:    "empty id",
			reqBody: validBody,
		},
		{
			name:    "id with white space",
			id:      "   ",
			reqBody: validBody,
		},
		{
			name:    "invalid id",
			id:      "invalid-id",
			reqBody: validBody,
		},
		{
			name: "empty title",
			id:   validID,
			reqBody: func() handler.UpdateBagRequestBody {
				b := validBody
				b.Title = ""
				return b
			}(),
		},
		{
			name: "title with white space",
			id:   validID,
			reqBody: func() handler.UpdateBagRequestBody {
				b := validBody
				b.Title = "   "
				return b
			}(),
		},
		{
			name: "title shorter than 3 characters",
			id:   validID,
			reqBody: func() handler.UpdateBagRequestBody {
				b := validBody
				b.Title = "ab"
				return b
			}(),
		},
		{
			name: "title longer than 50 characters",
			id:   validID,
			reqBody: func() handler.UpdateBagRequestBody {
				b := validBody
				b.Title = strings.Repeat("a", 51)
				return b
			}(),
		},
		{
			name: "description longer than 255 characters",
			id:   validID,
			reqBody: func() handler.UpdateBagRequestBody {
				b := validBody
				b.Description = strings.Repeat("a", 256)
				return b
			}(),
		},
		{
			name: "stock equal to 0",
			id:   validID,
			reqBody: func() handler.UpdateBagRequestBody {
				b := validBody
				b.Stock = 0
				return b
			}(),
		},
		{
			name: "negative stock",
			id:   validID,
			reqBody: func() handler.UpdateBagRequestBody {
				b := validBody
				b.Stock = -1
				return b
			}(),
		},
		{
			name: "originalPrice equal to 0",
			id:   validID,
			reqBody: func() handler.UpdateBagRequestBody {
				b := validBody
				b.OriginalPrice = 0
				return b
			}(),
		},
		{
			name: "negative originalPrice",
			id:   validID,
			reqBody: func() handler.UpdateBagRequestBody {
				b := validBody
				b.OriginalPrice = -1
				return b
			}(),
		},
		{
			name: "discount less than 30",
			id:   validID,
			reqBody: func() handler.UpdateBagRequestBody {
				b := validBody
				b.Discount = 29
				return b
			}(),
		},
		{
			name: "discount greater than 90",
			id:   validID,
			reqBody: func() handler.UpdateBagRequestBody {
				b := validBody
				b.Discount = 91
				return b
			}(),
		},
		{
			name: "empty items",
			id:   validID,
			reqBody: func() handler.UpdateBagRequestBody {
				b := validBody
				b.Items = []handler.UpdateBagItemRequestBody{}
				return b
			}(),
		},
		{
			name: "item with empty uuid",
			id:   validID,
			reqBody: func() handler.UpdateBagRequestBody {
				b := validBody
				b.Items = []handler.UpdateBagItemRequestBody{
					{Quantity: 2},
				}
				return b
			}(),
		},
		{
			name: "item with white space uuid",
			id:   validID,
			reqBody: func() handler.UpdateBagRequestBody {
				b := validBody
				b.Items = []handler.UpdateBagItemRequestBody{
					{ID: "   ", Quantity: 2},
				}
				return b
			}(),
		},
		{
			name: "item with invalid uuid",
			id:   validID,
			reqBody: func() handler.UpdateBagRequestBody {
				b := validBody
				b.Items = []handler.UpdateBagItemRequestBody{
					{ID: "invalid-id", Quantity: 1},
				}
				return b
			}(),
		},
		{
			name: "item with negative quantity",
			id:   validID,
			reqBody: func() handler.UpdateBagRequestBody {
				b := validBody
				b.Items = []handler.UpdateBagItemRequestBody{
					{ID: "550e8400-e29b-41d4-a716-446655440000", Quantity: -1},
				}
				return b
			}(),
		},
		{
			name: "item with quantity equal to 0",
			id:   validID,
			reqBody: func() handler.UpdateBagRequestBody {
				b := validBody
				b.Items = []handler.UpdateBagItemRequestBody{
					{ID: "550e8400-e29b-41d4-a716-446655440000", Quantity: 0},
				}
				return b
			}(),
		},
		{
			name: "end before start",
			id:   validID,
			reqBody: func() handler.UpdateBagRequestBody {
				b := validBody
				b.PickupWindow.End = now.Add(-time.Hour)
				return b
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clk := clock.NewFakeClock()
			bagRepo := repository.NewInMemRepo()
			productRepo := productRepository.NewInMemRepo()
			uc := usecase.NewUpdateBag(clk, bagRepo, productRepo)
			v := validation.NewValidator()
			h := handler.NewUpdateBag(v, uc)

			b, err := json.Marshal(tt.reqBody)
			require.NoError(t, err)

			req := newUpdateBagReq(tt.id, bytes.NewReader(b))
			rr := httptest.NewRecorder()

			h.Handle(rr, req)

			assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		})
	}
}
