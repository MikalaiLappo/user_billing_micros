package infrastructure

import (
	"context"
	"fmt"
	"sync"

	"github.com/MikalaiLappo/user_billing_micros/users/domain"
	"github.com/google/uuid"
)

type InMemoryUserRepository struct {
	mu       sync.RWMutex
	accounts map[string]*domain.User
}

func NewInMemoryBillingRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		accounts: make(map[string]*domain.User),
	}
}

func (r *InMemoryUserRepository) Save(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, hasOwnerAccount := r.accounts[user.ID.String()]; hasOwnerAccount {
		return fmt.Errorf("user with ID `%s` already exists", user.ID.String())
	}

	r.accounts[user.ID.String()] = user

	return nil
}

func (r *InMemoryUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, doesUserExist := r.accounts[id.String()]

	if !doesUserExist {
		return nil, fmt.Errorf("owner's `%s` billing account was not found", user.ID.String())
	}

	return user, nil
}
