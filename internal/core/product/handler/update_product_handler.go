package handler

import (
	"errors"
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/product/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/gin-gonic/gin"
)

type UpdateProductParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type UpdateProductRequest struct {
	Name        string `form:"name" binding:"required,notblank,min=3,max=100"`
	RemoveImage bool   `form:"removeImage"`
}

type UpdateProductHandler struct {
	updateProduct *usecase.UpdateProduct
}

func NewUpdateProductHandler(updateProduct *usecase.UpdateProduct) *UpdateProductHandler {
	return &UpdateProductHandler{
		updateProduct: updateProduct,
	}
}

func (h *UpdateProductHandler) Handle(c *gin.Context) {
	var params UpdateProductParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	var body UpdateProductRequest
	if err := c.ShouldBind(&body); err != nil {
		c.Error(err)
		return
	}

	var image usecase.ImageInput
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		if !errors.Is(err, http.ErrMissingFile) {
			c.Error(err)
			return
		}
	} else {
		defer file.Close()

		contentType, err := storage.GetContentType(file)
		if err != nil {
			c.Error(err)
			return
		}

		image = usecase.ImageInput{
			Content:     file,
			ContentType: contentType,
		}
	}

	if c.Request.MultipartForm != nil {
		defer c.Request.MultipartForm.RemoveAll()
	}

	storeID, ok := request.StoreID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := h.updateProduct.Execute(c.Request.Context(), usecase.UpdateProductInput{
		ID:          params.ID,
		StoreID:     storeID,
		Name:        body.Name,
		Image:       image,
		RemoveImage: body.RemoveImage,
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
