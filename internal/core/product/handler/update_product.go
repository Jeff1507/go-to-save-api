package handler

import (
	"errors"
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/product"
	"github.com/Jeff1507/go-to-save-api/internal/core/product/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/render"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/Jeff1507/go-to-save-api/internal/infra/validation"
	"github.com/go-chi/chi/v5"
)

type UpdateProductRequestBody struct {
	Name string
}

type UpdateProduct struct {
	uc *usecase.UpdateProduct
}

func NewUpdateProduct(uc *usecase.UpdateProduct) *UpdateProduct {
	return &UpdateProduct{
		uc: uc,
	}
}

func (h *UpdateProduct) Handle(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	defer r.Body.Close()

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		render.InvalidMultipartForm(w)
		return
	}
	defer r.MultipartForm.RemoveAll()

	id := chi.URLParam(r, "id")

	body := UpdateProductRequestBody{
		Name: r.FormValue("name"),
	}

	if errs := h.validate(id, body); errs.HasErrors() {
		render.ValidationError(w, errs)
		return
	}

	var imgInput *usecase.UpdateProductImageInput
	file, _, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		contentType, err := storage.GetContentType(file)
		if err != nil {
			render.InternalServerError(w)
			return
		}

		imgInput = &usecase.UpdateProductImageInput{
			Content:     file,
			ContentType: contentType,
		}
	}

	if err := h.uc.Execute(r.Context(), usecase.UpdateProductInput{
		ID:    id,
		Name:  body.Name,
		Image: imgInput,
	}); err != nil {
		if errors.Is(err, product.ErrInvalidImageFormat) {
			render.Error(w, http.StatusUnsupportedMediaType, "Formato de imagem não suportado")
			return
		}

		if errors.Is(err, product.ErrNotFound) {
			render.Error(w, http.StatusNotFound, "Produto não encontrado")
			return
		}

		if errors.Is(err, product.ErrAlreadyExists) {
			render.Error(w, http.StatusConflict, "Produto já cadastrado")
			return
		}

		render.InternalServerError(w)
		return
	}

	render.NoContent(w)
}

func (h UpdateProduct) validate(id string, body UpdateProductRequestBody) validation.Errors {
	errs := validation.NewErrors()

	validation.String(errs, "id", id).
		Required().
		UUID()

	validation.String(errs, "name", body.Name).
		Required().
		Min(3).
		Max(100)

	return errs
}
