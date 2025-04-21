package infrastructure

import (
	"context"
	"fmt"
	"sync"

	"github.com/MikalaiLappo/user_billing_micros/billing/application"
	"github.com/MikalaiLappo/user_billing_micros/billing/domain"
	"github.com/google/uuid"
)

type InMemoryBillingRepository struct {
	mu       sync.RWMutex
	accounts map[string]*domain.BillingAccount
}

func NewInMemoryBillingRepository() *InMemoryBillingRepository {
	return &InMemoryBillingRepository{
		accounts: make(map[string]*domain.BillingAccount),
	}
}

func (r *InMemoryBillingRepository) Save(ctx context.Context, account *domain.BillingAccount) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, hasOwnerAccount := r.accounts[account.OwnerID.String()]; hasOwnerAccount {
		return fmt.Errorf("owner `%s` already has a billing account", account.OwnerID.String())
	}

	r.accounts[account.OwnerID.String()] = account

	return nil
}

func (r *InMemoryBillingRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.BillingAccount, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	account, isAccountPresented := r.accounts[id.String()]

	if !isAccountPresented {
		return nil, fmt.Errorf("billing account of `%s` owner was not found", account.OwnerID.String())
	}

	return account, nil
}

func (r *InMemoryBillingRepository) Update(ctx context.Context, ownerID uuid.UUID, txFunc func(prev uint64) (uint64, error)) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	account, isAccountPresented := r.accounts[ownerID.String()]

	if !isAccountPresented {
		return fmt.Errorf("billing account of `%s` owner was not found", account.OwnerID.String())
	}

	newBalance, err := txFunc(account.Balance)
	fmt.Printf("new balance %d\n", newBalance)

	if err != nil {
		return err
	}

	r.accounts[ownerID.String()].Balance = newBalance

	return nil
}

var _ application.BillingRepository = &InMemoryBillingRepository{}
