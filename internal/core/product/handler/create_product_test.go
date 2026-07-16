package handler_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/product/handler"
	"github.com/Jeff1507/go-to-save-api/internal/core/product/repository"
	"github.com/Jeff1507/go-to-save-api/internal/core/product/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/Jeff1507/go-to-save-api/internal/infra/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateProduct_Success(t *testing.T) {
	clk := clock.NewFakeClock()
	st := storage.NewFakeStorage()
	repo := repository.NewInMemRepo()
	uc := usecase.NewCreateProduct(clk, st, repo)
	v := validation.NewValidator()
	h := handler.NewCreateProduct(v, uc)

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	err := w.WriteField("name", "Coxinha")
	require.NoError(t, err)

	part, err := w.CreateFormFile("image", "fake-img.jpg")
	require.NoError(t, err)

	_, err = part.Write([]byte("fake image content"))
	require.NoError(t, err)

	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/products", body)
	req.Header.Set("Content-Type", w.FormDataContentType())

	rr := httptest.NewRecorder()

	h.Handle(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestCreateProduct_ValidationError(t *testing.T) {
	tests := []struct {
		name    string
		reqBody handler.CreateProductRequest
	}{
		{
			name:    "empty name",
			reqBody: handler.CreateProductRequest{},
		},
		{
			name: "name with white space",
			reqBody: handler.CreateProductRequest{
				Name: "   ",
			},
		},
		{
			name: "name shorter than 3 characters",
			reqBody: handler.CreateProductRequest{
				Name: "ab",
			},
		},
		{
			name: "name longer than 100 characters",
			reqBody: handler.CreateProductRequest{
				Name: strings.Repeat("a", 101),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clk := clock.NewFakeClock()
			st := storage.NewFakeStorage()
			repo := repository.NewInMemRepo()
			uc := usecase.NewCreateProduct(clk, st, repo)
			v := validation.NewValidator()
			h := handler.NewCreateProduct(v, uc)

			body := &bytes.Buffer{}
			w := multipart.NewWriter(body)

			err := w.WriteField("name", tt.reqBody.Name)
			require.NoError(t, err)

			w.Close()

			req := httptest.NewRequest(http.MethodPost, "/products", body)
			req.Header.Set("Content-Type", w.FormDataContentType())

			rr := httptest.NewRecorder()

			h.Handle(rr, req)

			assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		})
	}
}
