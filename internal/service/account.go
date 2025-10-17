package service

import (
	"context"
	"errors"
	"time"

	"github.com/benrowe/nab-bank-api/internal/model"
)

// AccountService defines the interface for account operations
type AccountService interface {
	GetAllAccounts(ctx context.Context) ([]model.Account, error)
	GetAccountDetails(ctx context.Context, accountID string) (*model.AccountDetails, error)
}

// Service errors
var (
	ErrAccountNotFound      = errors.New("account not found")
	ErrServiceUnavailable   = errors.New("service unavailable")
	ErrAuthenticationFailed = errors.New("authentication failed")
)

// accountService implements AccountService
type accountService struct {
	nabClient NABClient
}

// NABClient defines the interface for interacting with NAB's website
type NABClient interface {
	GetAccounts(ctx context.Context) ([]model.Account, error)
	GetAccountTransactions(ctx context.Context, accountID string) ([]model.Transaction, error)
}

// NewAccountService creates a new account service
func NewAccountService(nabClient NABClient) AccountService {
	return &accountService{
		nabClient: nabClient,
	}
}

// GetAllAccounts retrieves all accounts from NAB
func (s *accountService) GetAllAccounts(ctx context.Context) ([]model.Account, error) {
	accounts, err := s.nabClient.GetAccounts(ctx)
	if err != nil {
		return nil, err
	}

	// Update last updated timestamp
	now := time.Now()
	for i := range accounts {
		accounts[i].LastUpdated = &now
	}

	return accounts, nil
}

// GetAccountDetails retrieves detailed account information including transactions
func (s *accountService) GetAccountDetails(ctx context.Context, accountID string) (*model.AccountDetails, error) {
	// First get all accounts to find the requested one
	accounts, err := s.nabClient.GetAccounts(ctx)
	if err != nil {
		return nil, err
	}

	var targetAccount *model.Account
	for _, account := range accounts {
		if account.ID == accountID {
			targetAccount = &account
			break
		}
	}

	if targetAccount == nil {
		return nil, ErrAccountNotFound
	}

	// Get transactions for this account
	transactions, err := s.nabClient.GetAccountTransactions(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// Update last updated timestamp
	now := time.Now()
	targetAccount.LastUpdated = &now

	accountDetails := &model.AccountDetails{
		Account:                  *targetAccount,
		Transactions:             transactions,
		RecentTransactionCount:   len(transactions),
	}

	return accountDetails, nil
}