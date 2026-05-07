package accounts

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type Account struct {
	ID   uuid.UUID
	Name string
}

type AccountManager struct {
	Accounts map[string]uuid.UUID
	mu       sync.Mutex
}

func NewAccountManager() *AccountManager {
	return &AccountManager{}
}

func (a *AccountManager) GetAccount(name string) (uuid.UUID, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	id, exists := a.Accounts[name]
	return id, exists
}

func (a *AccountManager) CreateAccount(name string) (*Account, error) {
	_, exists := a.GetAccount(name)
	if exists {
		return nil, fmt.Errorf("That name is already taken")
	}

	return &Account{
		ID:   uuid.New(),
		Name: name,
	}, nil
}

func (a *AccountManager) DeleteAccount(name string) (bool, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	_, exists := a.GetAccount(name)
	if !exists {
		return exists, fmt.Errorf("No such account exists")
	}

	delete(a.Accounts, name)

	_, exists = a.GetAccount(name)
	return exists, nil
}
