package handler

import "time"

type PickupWindowResponse struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}
