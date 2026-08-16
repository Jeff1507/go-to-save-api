package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/customer/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/gin-gonic/gin"
)

type UpdateProfileRequest struct {
	FirstName string `form:"firstName" binding:"required,notblank,min=2,max=50"`
	LastName  string `form:"lastName" binding:"required,notblank,min=2,max=50"`
}

type UpdateProfileHandler struct {
	updateProfile *usecase.UpdateProfile
}

func NewUpdateProfileHandler(updateProfile *usecase.UpdateProfile) *UpdateProfileHandler {
	return &UpdateProfileHandler{
		updateProfile: updateProfile,
	}
}

func (h *UpdateProfileHandler) Handle(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, storage.MaxUploadSize)

	if c.Request.MultipartForm != nil {
		defer c.Request.MultipartForm.RemoveAll()
	}

	var body UpdateProfileRequest
	if err := c.ShouldBind(&body); err != nil {
		c.Error(err)
		return
	}

	var avatar *usecase.UpdateProfileAvatarInput
	file, _, err := c.Request.FormFile("avatar")
	if err == nil {
		defer file.Close()

		contentType, err := storage.GetContentType(file)
		if err != nil {
			c.Error(err)
			return
		}

		avatar = &usecase.UpdateProfileAvatarInput{
			Content:     file,
			ContentType: contentType,
		}
	}

	id, ok := request.CustomerID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := h.updateProfile.Execute(c.Request.Context(), usecase.UpdateProfileInput{
		ID:        id,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Avatar:    avatar,
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
