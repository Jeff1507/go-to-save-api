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

type CreateBagItemRequest struct {
	ProductID string `json:"productId" validate:"required,uuid"`
	Quantity  int32  `json:"quantity" validate:"required,gt=0"`
}

type CreateBagPickupWindowRequest struct {
	Start time.Time `json:"start" validate:"required"`
	End   time.Time `json:"end" validate:"required,gtfield=Start"`
}

type CreateBagRequest struct {
	Title         string                       `json:"title" validate:"required,notblank,min=3,max=50"`
	Description   string                       `json:"description" validate:"omitempty,notblank,max=255"`
	Stock         int32                        `json:"stock" validate:"required,gt=0"`
	OriginalPrice int64                        `json:"originalPrice" validate:"required,gt=0"`
	Discount      int                          `json:"discount" validate:"required,gte=30,lte=90"`
	PickupWindow  CreateBagPickupWindowRequest `json:"pickupWindow" validate:"required"`
	Items         []CreateBagItemRequest       `json:"items" validate:"required,min=1,dive"`
}

type CreateBagResponse struct {
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

	var body CreateBagRequest
	if err := render.Bind(r, &body); err != nil {
		render.InvalidRequestBody(w)
		return
	}

	if errs := h.validateRequest(body); errs.HasErrors() {
		render.ValidationError(w, errs)
		return
	}

	// TODO: Pegar o storeID do ctx
	storeID := "b875eb25-fc40-48fc-8ed1-daff245deb94"

	pickupWindow := usecase.CreateBagPickupWindowInput{
		Start: body.PickupWindow.Start,
		End:   body.PickupWindow.End,
	}

	items := make([]usecase.CreateBagItemInput, 0, len(body.Items))
	for _, item := range body.Items {
		items = append(items, usecase.CreateBagItemInput{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	output, err := h.uc.Execute(r.Context(), usecase.CreateBagInput{
		StoreID:       storeID,
		Title:         body.Title,
		Description:   body.Description,
		Stock:         body.Stock,
		OriginalPrice: body.OriginalPrice,
		Discount:      body.Discount,
		PickupWindow:  pickupWindow,
		Items:         items,
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

	render.Success(w, http.StatusCreated, CreateBagResponse{
		ID: output.ID,
	})
}

func (h *CreateBag) validateRequest(body CreateBagRequest) validation.Errors {
	errs := validation.NewErrors()

	validation.Struct(h.v, errs, body)

	return errs
}
