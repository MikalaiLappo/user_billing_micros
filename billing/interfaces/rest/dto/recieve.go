package dto

import "github.com/go-playground/validator/v10"

type (
	RecieveRequest struct {
		OwnerID string `json:"owner_id" validate:"required"`
		Amount  uint64 `json:"amount" validate:"required"`
	}

	RecieveResponse struct {
		OwnerID    string `json:"owner_id"`
		NewBalance uint64 `json:"balance"`
	}
)

func (r *RecieveRequest) Validate() error {
	validator := validator.New()
	return validator.Struct(r)
}
