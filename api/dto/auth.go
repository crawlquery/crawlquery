package dto

import "encoding/json"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

type LoginResponse struct {
	Token string `json:"token"`
}

func NewLoginResponse(token string) *LoginResponse {
	return &LoginResponse{
		Token: token,
	}
}
