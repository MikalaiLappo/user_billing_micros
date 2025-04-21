package dto

import "github.com/go-playground/validator/v10"

type (
	PayRequest struct {
		OwnerID string `json:"owner_id" validate:"required"`
		Amount  uint64 `json:"amount" validate:"required"`
	}

	PayResponse struct {
		OwnerID    string `json:"owner_id"`
		NewBalance uint64 `json:"balance"`
	}
)

func (r *PayRequest) Validate() error {
	validator := validator.New()
	return validator.Struct(r)
}
