package handler

type ItemRequest struct {
	ProductID string `json:"productId" binding:"required,uuid"`
	Quantity  int32  `json:"quantity" binding:"required,gt=0"`
}
