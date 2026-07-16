package handler_test

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/product"
	"github.com/Jeff1507/go-to-save-api/internal/core/product/handler"
	"github.com/Jeff1507/go-to-save-api/internal/core/product/repository"
	"github.com/Jeff1507/go-to-save-api/internal/core/product/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/Jeff1507/go-to-save-api/internal/infra/validation"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newUpdateProductReq(id string, body io.Reader, contentType string) *http.Request {
	escapedID := url.PathEscape(id)
	req := httptest.NewRequest(http.MethodPut, "/products/"+escapedID, body)
	req.Header.Set("Content-Type", contentType)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)

	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestUpdateProduct_Success(t *testing.T) {
	clk := clock.NewFakeClock()
	st := storage.NewFakeStorage()
	repo := repository.NewInMemRepo()
	uc := usecase.NewUpdateProduct(clk, st, repo)
	v := validation.NewValidator()
	h := handler.NewUpdateProduct(v, uc)
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

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	err = w.WriteField("name", "Coxinha de frango")
	require.NoError(t, err)

	part, err := w.CreateFormFile("image", "coxinha-frango.png")
	require.NoError(t, err)

	_, err = part.Write([]byte("fake image content"))
	require.NoError(t, err)

	w.Close()

	req := newUpdateProductReq(prod.ID(), body, w.FormDataContentType())
	rr := httptest.NewRecorder()

	h.Handle(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestUpdateProduct_ValidationError(t *testing.T) {
	validID := "a3f8c2d1-1234-4abc-89ef-123456789abc"
	validName := "Coxinha de frango"

	tests := []struct {
		name    string
		id      string
		reqBody handler.UpdateProductRequest
	}{
		{
			name: "empty id",
			reqBody: handler.UpdateProductRequest{
				Name: validName,
			},
		},
		{
			name: "id with white space",
			id:   "   ",
			reqBody: handler.UpdateProductRequest{
				Name: validName,
			},
		},
		{
			name: "invalid id",
			id:   "invalid-id",
			reqBody: handler.UpdateProductRequest{
				Name: validName,
			},
		},
		{
			name:    "empty name",
			id:      validID,
			reqBody: handler.UpdateProductRequest{},
		},
		{
			name: "name with white space",
			id:   validID,
			reqBody: handler.UpdateProductRequest{
				Name: "   ",
			},
		},
		{
			name: "name shorter than 3 characters",
			id:   validID,
			reqBody: handler.UpdateProductRequest{
				Name: "ab",
			},
		},
		{
			name: "name longer than 100 characters",
			id:   validID,
			reqBody: handler.UpdateProductRequest{
				Name: strings.Repeat("a", 101),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clock := clock.NewFakeClock()
			st := storage.NewFakeStorage()
			repo := repository.NewInMemRepo()
			uc := usecase.NewUpdateProduct(clock, st, repo)
			v := validation.NewValidator()
			h := handler.NewUpdateProduct(v, uc)

			body := &bytes.Buffer{}
			w := multipart.NewWriter(body)

			err := w.WriteField("name", tt.reqBody.Name)
			require.NoError(t, err)

			w.Close()

			req := newUpdateProductReq(tt.id, body, w.FormDataContentType())
			rr := httptest.NewRecorder()

			h.Handle(rr, req)

			assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		})
	}
}
