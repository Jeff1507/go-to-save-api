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

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/handler"
	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/repository"
	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/validation"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newUpdatePickupTemplateReq(id string, body io.Reader) *http.Request {
	escapedID := url.PathEscape(id)
	req := httptest.NewRequest(http.MethodPut, "/pickup-templates/"+escapedID, body)
	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)

	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestUpdatePickupTemplate_Success(t *testing.T) {
	repo := repository.NewInMemRepo()
	uc := usecase.NewUpdatePickupTemplate(repo)
	v := validation.NewValidator()
	h := handler.NewUpdatePickupTemplate(v, uc)

	start, err := pickup.NewTime(20, 0)
	require.NoError(t, err)

	end, err := pickup.NewTime(0, 0)
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

	body := handler.UpdatePickupTemplateRequestBody{
		Name:  "Fechamento",
		Start: "21:00",
		End:   "01:00",
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	req := newUpdatePickupTemplateReq(templ.ID(), bytes.NewReader(b))
	rr := httptest.NewRecorder()

	h.Handle(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestUpdatePickupTemplate_ValidationError(t *testing.T) {
	validID := "a3f8c2d1-1234-4abc-89ef-123456789abc"
	validName := "Fechamento"
	validStart := "21:00"
	validEnd := "01:00"

	tests := []struct {
		name    string
		id      string
		reqBody handler.UpdatePickupTemplateRequestBody
	}{
		{
			name: "empty id",
			reqBody: handler.UpdatePickupTemplateRequestBody{
				Name:  validName,
				Start: validStart,
				End:   validEnd,
			},
		},
		{
			name: "id with white space",
			id:   "   ",
			reqBody: handler.UpdatePickupTemplateRequestBody{
				Name:  validName,
				Start: validStart,
				End:   validEnd,
			},
		},
		{
			name: "invalid id",
			id:   "invalid-id",
			reqBody: handler.UpdatePickupTemplateRequestBody{
				Name:  validName,
				Start: validStart,
				End:   validEnd,
			},
		},
		{
			name: "empty name",
			id:   validID,
			reqBody: handler.UpdatePickupTemplateRequestBody{
				Start: validStart,
				End:   validEnd,
			},
		},
		{
			name: "name with white space",
			id:   validID,
			reqBody: handler.UpdatePickupTemplateRequestBody{
				Name:  "   ",
				Start: validStart,
				End:   validEnd,
			},
		},
		{
			name: "name shorter than 3 characters",
			id:   validID,
			reqBody: handler.UpdatePickupTemplateRequestBody{
				Name:  "ab",
				Start: validStart,
				End:   validEnd,
			},
		},
		{
			name: "name longer than 50 characters",
			id:   validID,
			reqBody: handler.UpdatePickupTemplateRequestBody{
				Name:  strings.Repeat("a", 51),
				Start: validStart,
				End:   validEnd,
			},
		},
		{
			name: "empty start time",
			id:   validID,
			reqBody: handler.UpdatePickupTemplateRequestBody{
				Name: validName,
				End:  validEnd,
			},
		},
		{
			name: "start with white space",
			id:   validID,
			reqBody: handler.UpdatePickupTemplateRequestBody{
				Name:  validName,
				Start: "   ",
				End:   validEnd,
			},
		},
		{
			name: "invalid start time",
			id:   validID,
			reqBody: handler.UpdatePickupTemplateRequestBody{
				Name:  validName,
				Start: "xdd:xdd",
				End:   validEnd,
			},
		},
		{
			name: "start hours greater than 23",
			id:   validID,
			reqBody: handler.UpdatePickupTemplateRequestBody{
				Name:  validName,
				Start: "24:00",
				End:   validEnd,
			},
		},
		{
			name: "start hours less than 0",
			id:   validID,
			reqBody: handler.UpdatePickupTemplateRequestBody{
				Name:  validName,
				Start: "-00:01",
				End:   validEnd,
			},
		},
		{
			name: "start minutes greater than 59",
			id:   validID,
			reqBody: handler.UpdatePickupTemplateRequestBody{
				Name:  validName,
				Start: "23:60",
				End:   validEnd,
			},
		},
		{
			name: "start minutes less than 0",
			id:   validID,
			reqBody: handler.UpdatePickupTemplateRequestBody{
				Name:  validName,
				Start: "00:-01",
				End:   validEnd,
			},
		},
		{
			name: "empty end time",
			id:   validID,
			reqBody: handler.UpdatePickupTemplateRequestBody{
				Name:  validName,
				Start: validStart,
			},
		},
		{
			name: "end with white space",
			id:   validID,
			reqBody: handler.UpdatePickupTemplateRequestBody{
				Name:  validName,
				Start: validStart,
				End:   "   ",
			},
		},
		{
			name: "invalid end time",
			id:   validID,
			reqBody: handler.UpdatePickupTemplateRequestBody{
				Name:  validName,
				Start: validStart,
				End:   "xdd:xdd",
			},
		},
		{
			name: "end hours greater than 23",
			id:   validID,
			reqBody: handler.UpdatePickupTemplateRequestBody{
				Name:  validName,
				Start: validStart,
				End:   "24:00",
			},
		},
		{
			name: "end hours less than 0",
			id:   validID,
			reqBody: handler.UpdatePickupTemplateRequestBody{
				Name:  validName,
				Start: validStart,
				End:   "-00:01",
			},
		},
		{
			name: "end minutes greater than 59",
			id:   validID,
			reqBody: handler.UpdatePickupTemplateRequestBody{
				Name:  validName,
				Start: validStart,
				End:   "23:60",
			},
		},
		{
			name: "end minutes less than 0",
			id:   validID,
			reqBody: handler.UpdatePickupTemplateRequestBody{
				Name:  validName,
				Start: validStart,
				End:   "00:-01",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repository.NewInMemRepo()
			uc := usecase.NewUpdatePickupTemplate(repo)
			v := validation.NewValidator()
			h := handler.NewUpdatePickupTemplate(v, uc)

			body := handler.UpdatePickupTemplateRequestBody{
				Name:  tt.reqBody.Name,
				Start: tt.reqBody.Start,
				End:   tt.reqBody.End,
			}

			b, err := json.Marshal(body)
			require.NoError(t, err)

			req := newUpdatePickupTemplateReq(tt.id, bytes.NewReader(b))
			rr := httptest.NewRecorder()

			h.Handle(rr, req)

			assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		})
	}
}
