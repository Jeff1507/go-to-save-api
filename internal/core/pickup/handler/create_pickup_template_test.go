package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/handler"
	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/repository"
	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePickupTemplate_Success(t *testing.T) {
	repo := repository.NewInMemRepo()
	uc := usecase.NewCreatePickupTemplate(repo)
	v := validation.NewValidator()
	h := handler.NewCreatePickupTemplate(v, uc)

	body := handler.CreatePickupTemplateRequest{
		Name:  "Domingo a noite",
		Start: "20:00",
		End:   "22:00",
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/pickup-templates", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	h.Handle(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestCreatePickupTemplate_ValidationError(t *testing.T) {
	validName := "Fechamento"
	validStart := "22:00"
	validEnd := "00:00"

	tests := []struct {
		name    string
		reqBody handler.CreatePickupTemplateRequest
	}{
		{
			name: "empty name",
			reqBody: handler.CreatePickupTemplateRequest{
				Start: validStart,
				End:   validEnd,
			},
		},
		{
			name: "name with white space",
			reqBody: handler.CreatePickupTemplateRequest{
				Name:  "   ",
				Start: validStart,
				End:   validEnd,
			},
		},
		{
			name: "name shorter than 3 characters",
			reqBody: handler.CreatePickupTemplateRequest{
				Name:  "ab",
				Start: validStart,
				End:   validEnd,
			},
		},
		{
			name: "name longer than 50 characters",
			reqBody: handler.CreatePickupTemplateRequest{
				Name:  strings.Repeat("a", 51),
				Start: validStart,
				End:   validEnd,
			},
		},
		{
			name: "empty start time",
			reqBody: handler.CreatePickupTemplateRequest{
				Name: validName,
				End:  validEnd,
			},
		},
		{
			name: "start with white space",
			reqBody: handler.CreatePickupTemplateRequest{
				Name:  validName,
				Start: "   ",
				End:   validEnd,
			},
		},
		{
			name: "invalid start time",
			reqBody: handler.CreatePickupTemplateRequest{
				Name:  validName,
				Start: "xdd:xdd",
				End:   validEnd,
			},
		},
		{
			name: "start hours greater than 23",
			reqBody: handler.CreatePickupTemplateRequest{
				Name:  validName,
				Start: "24:00",
				End:   validEnd,
			},
		},
		{
			name: "start hours less than 0",
			reqBody: handler.CreatePickupTemplateRequest{
				Name:  validName,
				Start: "-00:01",
				End:   validEnd,
			},
		},
		{
			name: "start minutes greater than 59",
			reqBody: handler.CreatePickupTemplateRequest{
				Name:  validName,
				Start: "23:60",
				End:   validEnd,
			},
		},
		{
			name: "start minutes less than 0",
			reqBody: handler.CreatePickupTemplateRequest{
				Name:  validName,
				Start: "00:-01",
				End:   validEnd,
			},
		},
		{
			name: "empty end time",
			reqBody: handler.CreatePickupTemplateRequest{
				Name:  validName,
				Start: validStart,
			},
		},
		{
			name: "end with white space",
			reqBody: handler.CreatePickupTemplateRequest{
				Name:  validName,
				Start: validStart,
				End:   "   ",
			},
		},
		{
			name: "invalid end time",
			reqBody: handler.CreatePickupTemplateRequest{
				Name:  validName,
				Start: validStart,
				End:   "xdd:xdd",
			},
		},
		{
			name: "end hours greater than 23",
			reqBody: handler.CreatePickupTemplateRequest{
				Name:  validName,
				Start: validStart,
				End:   "24:00",
			},
		},
		{
			name: "end hours less than 0",
			reqBody: handler.CreatePickupTemplateRequest{
				Name:  validName,
				Start: validStart,
				End:   "-00:01",
			},
		},
		{
			name: "end minutes greater than 59",
			reqBody: handler.CreatePickupTemplateRequest{
				Name:  validName,
				Start: validStart,
				End:   "23:60",
			},
		},
		{
			name: "end minutes less than 0",
			reqBody: handler.CreatePickupTemplateRequest{
				Name:  validName,
				Start: validStart,
				End:   "00:-01",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repository.NewInMemRepo()
			uc := usecase.NewCreatePickupTemplate(repo)
			v := validation.NewValidator()
			h := handler.NewCreatePickupTemplate(v, uc)

			body := handler.CreatePickupTemplateRequest{
				Name:  tt.reqBody.Name,
				Start: tt.reqBody.Start,
				End:   tt.reqBody.End,
			}

			b, err := json.Marshal(body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/pickup-templates", bytes.NewReader(b))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			h.Handle(rr, req)

			assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		})
	}
}
