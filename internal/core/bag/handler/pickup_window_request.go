package handler

type PickupWindowRequest struct {
	Start string `json:"start" binding:"required,time"`
	End   string `json:"end" binding:"required,time"`
}
