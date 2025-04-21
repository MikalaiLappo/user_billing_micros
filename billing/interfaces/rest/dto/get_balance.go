package dto

import "github.com/go-playground/validator/v10"

type (
	GetBalanceRequest struct {
		OwnerID string `json:"owner_id" validate:"required"`
	}

	GetBalanceResponse struct {
		OwnerID string `json:"owner_id"`
		Balance uint64 `json:"balance"`
	}
)

func (r *GetBalanceRequest) Validate() error {
	validator := validator.New()
	return validator.Struct(r)
}
