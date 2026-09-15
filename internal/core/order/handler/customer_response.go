package handler

type CustomerResponse struct {
	ID        string  `json:"id"`
	FullName  string  `json:"fullname"`
	AvatarURL *string `json:"avatarUrl"`
}
