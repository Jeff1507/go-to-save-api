package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/bag"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/core/product"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/render"
	"github.com/Jeff1507/go-to-save-api/internal/infra/validation"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type UpdateBagItemRequest struct {
	ProductID string `json:"productId" validate:"required,uuid"`
	Quantity  int32  `json:"quantity" validate:"required,gt=0"`
}

type UpdateBagPickupWindowRequest struct {
	Start time.Time `json:"start" validate:"required"`
	End   time.Time `json:"end" validate:"required,gtfield=Start"`
}

type UpdateBagRequest struct {
	Title         string                       `json:"title" validate:"required,notblank,min=3,max=50"`
	Description   string                       `json:"description" validate:"omitempty,notblank,max=255"`
	Stock         int32                        `json:"stock" validate:"required,gt=0"`
	OriginalPrice int64                        `json:"originalPrice" validate:"required,gt=0"`
	Discount      int                          `json:"discount" validate:"required,gte=30,lte=90"`
	Items         []UpdateBagItemRequest       `json:"items" validate:"required,min=1,dive"`
	PickupWindow  UpdateBagPickupWindowRequest `json:"pickupWindow" validate:"required"`
}

type UpdateBag struct {
	v  *validator.Validate
	uc *usecase.UpdateBag
}

func NewUpdateBag(v *validator.Validate, uc *usecase.UpdateBag) *UpdateBag {
	return &UpdateBag{
		v:  v,
		uc: uc,
	}
}

func (h *UpdateBag) Handle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id := chi.URLParam(r, "id")

	var body UpdateBagRequest
	if err := render.Bind(r, &body); err != nil {
		render.InvalidRequestBody(w)
		return
	}

	if errs := h.validateRequest(id, body); errs.HasErrors() {
		render.ValidationError(w, errs)
		return
	}

	pickupWindow := usecase.UpdateBagPickupWindowInput{
		Start: body.PickupWindow.Start,
		End:   body.PickupWindow.End,
	}

	items := make([]usecase.UpdateBagItemInput, 0, len(body.Items))
	for _, item := range body.Items {
		items = append(items, usecase.UpdateBagItemInput{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	if err := h.uc.Execute(r.Context(), usecase.UpdateBagInput{
		ID:            id,
		Title:         body.Title,
		Description:   body.Description,
		Stock:         body.Stock,
		OriginalPrice: body.OriginalPrice,
		Discount:      body.Discount,
		PickupWindow:  pickupWindow,
		Items:         items,
	}); err != nil {
		if errors.Is(err, bag.ErrNotFound) {
			render.Error(w, http.StatusNotFound, "Sacola não encontrada")
			return
		}

		if errors.Is(err, bag.ErrClosed) {
			render.Error(w, http.StatusForbidden, "Sacola encerrada")
			return
		}

		if errors.Is(err, bag.ErrExpired) {
			render.Error(w, http.StatusBadRequest, "Sacola expirada")
			return
		}

		if errors.Is(err, bag.ErrInvalidPickupWindowInterval) {
			render.Error(w, http.StatusUnprocessableEntity, "Janela de retirada inválida")
			return
		}

		if errors.Is(err, product.ErrNotFound) {
			render.Error(w, http.StatusNotFound, "Produto não encontrado")
			return
		}

		render.InternalServerError(w)
		return
	}

	render.NoContent(w)
}

func (h *UpdateBag) validateRequest(id string, body UpdateBagRequest) validation.Errors {
	errs := validation.NewErrors()

	validation.UUID(h.v, errs, id, "id")

	validation.Struct(h.v, errs, body)

	return errs
}
