package handler

import (
	"errors"
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/product"
	"github.com/Jeff1507/go-to-save-api/internal/core/product/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/render"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/Jeff1507/go-to-save-api/internal/infra/validation"
)

type CreateProductRequestBody struct {
	Name string
}

type CreateProductResponseBody struct {
	ID string `json:"id"`
}

type CreateProduct struct {
	uc *usecase.CreateProduct
}

func NewCreateProduct(uc *usecase.CreateProduct) *CreateProduct {
	return &CreateProduct{
		uc: uc,
	}
}

func (h *CreateProduct) Handle(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	defer r.Body.Close()

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		render.InvalidMultipartForm(w)
		return
	}
	defer r.MultipartForm.RemoveAll()

	body := CreateProductRequestBody{
		Name: r.FormValue("name"),
	}

	if errs := h.validate(body); errs.HasErrors() {
		render.ValidationError(w, errs)
		return
	}

	var imgInput *usecase.CreateProductImageInput
	file, _, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		contentType, err := storage.GetContentType(file)
		if err != nil {
			render.InternalServerError(w)
			return
		}

		imgInput = &usecase.CreateProductImageInput{
			Content:     file,
			ContentType: contentType,
		}
	}

	// TODO: Pegar o storeID do ctx
	storeID := "store-id"

	output, err := h.uc.Execute(r.Context(), usecase.CreateProductInput{
		StoreID: storeID,
		Name:    body.Name,
		Image:   imgInput,
	})
	if err != nil {
		if errors.Is(err, product.ErrInvalidImageFormat) {
			render.Error(w, http.StatusUnsupportedMediaType, "Formato de imagem não suportado")
			return
		}

		if errors.Is(err, product.ErrAlreadyExists) {
			render.Error(w, http.StatusConflict, "Produto já cadastrado")
			return
		}

		render.InternalServerError(w)
		return
	}

	render.Success(w, http.StatusCreated, CreateProductResponseBody{
		ID: output.ID,
	})
}

func (h CreateProduct) validate(body CreateProductRequestBody) validation.Errors {
	errs := validation.NewErrors()

	validation.String(errs, "name", body.Name).
		Required().
		Min(3).
		Max(100)

	return errs
}
