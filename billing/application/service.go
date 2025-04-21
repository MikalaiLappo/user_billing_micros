package application

import (
	"context"
	"fmt"

	"github.com/MikalaiLappo/user_billing_micros/billing/domain"
	"github.com/google/uuid"
)

type BillingService struct {
	repo BillingRepository
}

func NewBillingService(repo BillingRepository) *BillingService {
	return &BillingService{
		repo: repo,
	}
}

type BillingRepository interface {
	Save(ctx context.Context, account *domain.BillingAccount) error
	FindByID(ctx context.Context, ownerID uuid.UUID) (*domain.BillingAccount, error)
	Update(ctx context.Context, ownerID uuid.UUID, txFunc func(balance uint64) (uint64, error)) error
}

func (s *BillingService) GetBillingBalance(ctx context.Context, id uuid.UUID) (uint64, error) {
	account, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return 0, err
	}

	return account.Balance, nil
}

func (s *BillingService) Pay(ctx context.Context, ownerID uuid.UUID, amount uint64) error {
	return s.repo.Update(ctx, ownerID, func(balance uint64) (uint64, error) {
		if balance < amount {
			return balance, fmt.Errorf("insufficient funds: (balance) %d (pay request) %d", balance, amount)
		}
		return balance - amount, nil
	})
}

func (s *BillingService) RecieveMoney(ctx context.Context, ownerID uuid.UUID, amount uint64) error {
	return s.repo.Update(ctx, ownerID, func(balance uint64) (uint64, error) {
		fmt.Printf("recieve balance account: %d | update amount: %d\n", balance, amount)
		return balance + amount, nil
	})
}
