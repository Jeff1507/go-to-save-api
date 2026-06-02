package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag/handler"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag/repository"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/core/product"
	productRepository "github.com/Jeff1507/go-to-save-api/internal/core/product/repository"
	"github.com/Jeff1507/go-to-save-api/internal/infra/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateBag_Success(t *testing.T) {
	clk := clock.NewFakeClock()
	bagRepo := repository.NewInMemRepo()
	productRepo := productRepository.NewInMemRepo()
	uc := usecase.NewCreateBag(clk, bagRepo, productRepo)
	v := validation.NewValidator()
	h := handler.NewCreateBag(v, uc)

	now := clk.Now()

	prod, err := product.New(
		"a3f8c2d1-1234-4abc-89ef-123456789abc",
		"store-id",
		"Coxinha",
		"coxinha.png",
		now,
	)
	require.NoError(t, err)

	err = productRepo.Save(context.Background(), prod)
	require.NoError(t, err)

	pickupStart := now.Add(time.Hour * 12)
	pickupEnd := pickupStart.Add(time.Hour)

	body := handler.CreateBagRequestBody{
		Title:         "Sacola mista",
		Description:   "Sacola com doces e salgados",
		Stock:         1,
		OriginalPrice: 1000,
		Discount:      50,
		PickupWindow: handler.CreateBagPickupWindowRequestBody{
			Start: pickupStart,
			End:   pickupEnd,
		},
		Items: []handler.CreateBagItemRequestBody{
			{ID: prod.ID(), Quantity: 2},
		},
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/bags", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	h.Handle(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestCreateBag_ValidationError(t *testing.T) {
	v := validation.NewValidator()

	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	validBody := handler.CreateBagRequestBody{
		Title:         "Sacola teste",
		Description:   "",
		Stock:         10,
		OriginalPrice: 100,
		Discount:      30,
		Items: []handler.CreateBagItemRequestBody{
			{ID: "550e8400-e29b-41d4-a716-446655440000", Quantity: 1},
		},
		PickupWindow: handler.CreateBagPickupWindowRequestBody{
			Start: now,
			End:   now.Add(time.Hour),
		},
	}

	tests := []struct {
		name   string
		body   handler.CreateBagRequestBody
		fields []string
	}{
		{
			name:   "valid request body",
			body:   validBody,
			fields: nil,
		},
		{
			name: "empty title",
			body: func() handler.CreateBagRequestBody {
				b := validBody
				b.Title = ""
				return b
			}(),
			fields: []string{"title"},
		},
		{
			name: "title with white space",
			body: func() handler.CreateBagRequestBody {
				b := validBody
				b.Title = "   "
				return b
			}(),
			fields: []string{"title"},
		},
		{
			name: "title shorter than 3 characters",
			body: func() handler.CreateBagRequestBody {
				b := validBody
				b.Title = "ab"
				return b
			}(),
			fields: []string{"title"},
		},
		{
			name: "title longer than 50 characters",
			body: func() handler.CreateBagRequestBody {
				b := validBody
				b.Title = strings.Repeat("a", 51)
				return b
			}(),
			fields: []string{"title"},
		},
		{
			name: "description longer than 255 characters",
			body: func() handler.CreateBagRequestBody {
				b := validBody
				b.Description = strings.Repeat("a", 256)
				return b
			}(),
			fields: []string{"description"},
		},
		{
			name: "stock equal to 0",
			body: func() handler.CreateBagRequestBody {
				b := validBody
				b.Stock = 0
				return b
			}(),
			fields: []string{"stock"},
		},
		{
			name: "negative stock",
			body: func() handler.CreateBagRequestBody {
				b := validBody
				b.Stock = -1
				return b
			}(),
			fields: []string{"stock"},
		},
		{
			name: "originalPrice equal to 0",
			body: func() handler.CreateBagRequestBody {
				b := validBody
				b.OriginalPrice = 0
				return b
			}(),
			fields: []string{"originalPrice"},
		},
		{
			name: "negative originalPrice",
			body: func() handler.CreateBagRequestBody {
				b := validBody
				b.OriginalPrice = -1
				return b
			}(),
			fields: []string{"originalPrice"},
		},
		{
			name: "discount less than 30",
			body: func() handler.CreateBagRequestBody {
				b := validBody
				b.Discount = 29
				return b
			}(),
			fields: []string{"discount"},
		},
		{
			name: "discount greater than 90",
			body: func() handler.CreateBagRequestBody {
				b := validBody
				b.Discount = 91
				return b
			}(),
			fields: []string{"discount"},
		},
		{
			name: "item with empty uuid",
			body: func() handler.CreateBagRequestBody {
				b := validBody
				b.Items = []handler.CreateBagItemRequestBody{
					{Quantity: 1},
				}
				return b
			}(),
			fields: []string{"id"},
		},
		{
			name: "item with white space uuid",
			body: func() handler.CreateBagRequestBody {
				b := validBody
				b.Items = []handler.CreateBagItemRequestBody{
					{ID: "   ", Quantity: 1},
				}
				return b
			}(),
			fields: []string{"id"},
		},
		{
			name: "item with invalid uuid",
			body: func() handler.CreateBagRequestBody {
				b := validBody
				b.Items = []handler.CreateBagItemRequestBody{
					{ID: "invalid-id", Quantity: 1},
				}
				return b
			}(),
			fields: []string{"id"},
		},
		{
			name: "item with negative quantity",
			body: func() handler.CreateBagRequestBody {
				b := validBody
				b.Items = []handler.CreateBagItemRequestBody{
					{ID: "550e8400-e29b-41d4-a716-446655440000", Quantity: -1},
				}
				return b
			}(),
			fields: []string{"quantity"},
		},
		{
			name: "item with quantity equal to 0",
			body: func() handler.CreateBagRequestBody {
				b := validBody
				b.Items = []handler.CreateBagItemRequestBody{
					{ID: "550e8400-e29b-41d4-a716-446655440000", Quantity: 0},
				}
				return b
			}(),
			fields: []string{"quantity"},
		},
		{
			name: "end before start",
			body: func() handler.CreateBagRequestBody {
				b := validBody
				b.PickupWindow.End = now.Add(-time.Hour)
				return b
			}(),
			fields: []string{"end"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := validation.NewErrors()
			validation.Struct(v, errs, tt.body)

			if tt.fields == nil {
				if errs.HasErrors() {
					assert.NoError(t, errs)
				}
				return
			}

			for _, field := range tt.fields {
				if !errs.Has(field) {
					assert.Error(t, errs, field, errs)
				}
			}
		})
	}
}
