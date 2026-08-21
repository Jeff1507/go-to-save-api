package handler

import (
	"errors"
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/product/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/gin-gonic/gin"
)

type CreateProductRequest struct {
	Name string `form:"name" binding:"required,notblank,min=3,max=100"`
}

type CreateProductResponse struct {
	ID string `json:"id"`
}

type CreateProductHandler struct {
	createProduct *usecase.CreateProduct
}

func NewCreateProductHandler(createProduct *usecase.CreateProduct) *CreateProductHandler {
	return &CreateProductHandler{
		createProduct: createProduct,
	}
}

func (h *CreateProductHandler) Handle(c *gin.Context) {
	var body CreateProductRequest
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

	output, err := h.createProduct.Execute(c.Request.Context(), usecase.CreateProductInput{
		StoreID: storeID,
		Name:    body.Name,
		Image:   image,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": CreateProductResponse{
			ID: output.ID,
		},
	})
}
