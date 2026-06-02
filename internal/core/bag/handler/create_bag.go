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
	"github.com/go-playground/validator/v10"
)

type CreateBagItemRequestBody struct {
	ID       string `json:"id" validate:"required,uuid"`
	Quantity int32  `json:"quantity" validate:"required,gt=0"`
}

type CreateBagPickupWindowRequestBody struct {
	Start time.Time `json:"start" validate:"required"`
	End   time.Time `json:"end" validate:"required,gtfield=Start"`
}

type CreateBagRequestBody struct {
	Title         string                           `json:"title" validate:"required,notblank,min=3,max=50"`
	Description   string                           `json:"description" validate:"omitempty,notblank,max=255"`
	Stock         int32                            `json:"stock" validate:"required,gt=0"`
	OriginalPrice int64                            `json:"originalPrice" validate:"required,gt=0"`
	Discount      int                              `json:"discount" validate:"required,gte=30,lte=90"`
	Items         []CreateBagItemRequestBody       `json:"items" validate:"required,min=1,dive"`
	PickupWindow  CreateBagPickupWindowRequestBody `json:"pickupWindow" validate:"required"`
}

type CreateBagResponseBody struct {
	ID string `json:"id"`
}

type CreateBag struct {
	v  *validator.Validate
	uc *usecase.CreateBag
}

func NewCreateBag(v *validator.Validate, uc *usecase.CreateBag) *CreateBag {
	return &CreateBag{
		v:  v,
		uc: uc,
	}
}

func (h *CreateBag) Handle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var body CreateBagRequestBody
	if err := render.Bind(r, &body); err != nil {
		render.InvalidRequestBody(w)
		return
	}

	if errs := h.validateRequest(body); errs.HasErrors() {
		render.ValidationError(w, errs)
		return
	}

	// TODO: Pegar o storeID do ctx
	storeID := "store-id"

	inputItems := make([]usecase.CreateBagItemInput, 0, len(body.Items))
	for _, item := range body.Items {
		inputItems = append(inputItems, usecase.CreateBagItemInput{
			ID:       item.ID,
			Quantity: item.Quantity,
		})
	}

	output, err := h.uc.Execute(r.Context(), usecase.CreateBagInput{
		StoreID:       storeID,
		Title:         body.Title,
		Description:   body.Description,
		Stock:         body.Stock,
		OriginalPrice: body.OriginalPrice,
		Discount:      body.Discount,
		Items:         inputItems,
		PickupWindow: usecase.CreateBagPickupWindowInput{
			Start: body.PickupWindow.Start,
			End:   body.PickupWindow.End,
		},
	})
	if err != nil {
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

	render.Success(w, http.StatusCreated, CreateBagResponseBody{
		ID: output.ID,
	})
}

func (h *CreateBag) validateRequest(body CreateBagRequestBody) validation.Errors {
	errs := validation.NewErrors()

	validation.Struct(h.v, errs, body)

	return errs
}
