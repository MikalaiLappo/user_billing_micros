package domain

import "github.com/google/uuid"

type BillingAccount struct {
	OwnerID uuid.UUID
	Balance uint64
}

func NewBillingAccount() *BillingAccount {
	return &BillingAccount{
		OwnerID: uuid.New(),
		Balance: 0,
	}
}
