package services

import (
	"github.com/SwanHtetAungPhyo/common/model"
)

// AccountService defines methods related to account management.
type AccountService struct{}

// NewAccount creates a new account.
func (s *AccountService) NewAccount() (*model.Wallet, string, error) {
	return model.NewAccount()
}
