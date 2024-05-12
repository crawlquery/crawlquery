package dto

import (
	"crawlquery/api/domain"
	"encoding/json"
	"time"
)

type CreateAccountRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *CreateAccountRequest) ToJson() ([]byte, error) {
	return json.Marshal(r)
}

type CreateAccountResponse struct {
	Account struct {
		ID        string    `json:"id"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"account"`
}

func NewCreateAccountResponse(a *domain.Account) *CreateAccountResponse {
	res := &CreateAccountResponse{}

	res.Account.ID = a.ID
	res.Account.Email = a.Email
	res.Account.CreatedAt = a.CreatedAt

	return res
}
