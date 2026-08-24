package handler

type AddressRequest struct {
	PostalCode   string   `json:"postalCode" binding:"required,notblank,numeric,len=8"`
	Street       string   `json:"street" binding:"required,notblank,max=150"`
	Number       string   `json:"number" binding:"required,notblank,max=20"`
	Neighborhood string   `json:"neighborhood" binding:"required,notblank,max=100"`
	Complement   string   `json:"complement" binding:"omitempty,notblank,max=255"`
	City         string   `json:"city" binding:"required,notblank,max=100"`
	State        string   `json:"state" binding:"required,notblank,alpha,len=2"`
	Latitude     *float64 `json:"latitude" binding:"required,latitude"`
	Longitude    *float64 `json:"longitude" binding:"required,longitude"`
}
