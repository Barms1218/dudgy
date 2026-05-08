package accounts

import (
	"fmt"
	t "github.com/Barms1218/dudgy/internal/types"
	"sync"

	"github.com/google/uuid"
)

type AccountManager struct {
	Accounts map[uuid.UUID]*t.Account
	mu       sync.Mutex
}

func NewAccountManager() *AccountManager {
	return &AccountManager{}
}

func (a *AccountManager) GetAccount(id uuid.UUID) (*t.Account, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	account, exists := a.Accounts[id]
	return account, exists
}

func (a *AccountManager) CreateAccount(id uuid.UUID, name string) (*t.Account, error) {
	_, exists := a.GetAccount(id)
	if exists {
		return nil, fmt.Errorf("That name is already taken")
	}

	return &t.Account{
		ID:   uuid.New(),
		Name: name,
	}, nil
}

func (a *AccountManager) DeleteAccount(id uuid.UUID) (bool, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	_, exists := a.GetAccount(id)
	if !exists {
		return exists, fmt.Errorf("No such account exists")
	}

	delete(a.Accounts, id)

	_, exists = a.GetAccount(id)
	return exists, nil
}
