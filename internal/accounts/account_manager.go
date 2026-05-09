package accounts

import (
	"fmt"
	t "github.com/Barms1218/dudgy/internal/types"
	"sync"
)

type AccountManager struct {
	Accounts map[string]*t.Account
	mu       sync.Mutex
}

func NewAccountManager() *AccountManager {
	return &AccountManager{}
}

func (a *AccountManager) GetAccount(id string) *t.Account {
	a.mu.Lock()
	defer a.mu.Unlock()
	account, exists := a.Accounts[id]
	if !exists {
		return nil
	}
	return account
}

func (a *AccountManager) SetAccountName(id string, name string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	account, exists := a.Accounts[id]
	if !exists {
		return fmt.Errorf("No account with that id exists")
	}

	account.Name = name

	return nil
}

func (a *AccountManager) GetOrCreateAccount(id string) *t.Account {
	a.mu.Lock()
	defer a.mu.Unlock()
	account, exists := a.Accounts[id]
	if exists {
		return account
	}

	newAccount := &t.Account{
		ID: id,
	}
	a.Accounts[id] = newAccount
	return newAccount
}

func (a *AccountManager) DeleteAccount(id string) (bool, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	account := a.GetAccount(id)
	if account == nil {
		return false, fmt.Errorf("No such account exists")
	}

	delete(a.Accounts, id)

	account = a.GetAccount(id)
	return false, nil
}
