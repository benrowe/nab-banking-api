package service

import (
	"context"
	"time"

	"github.com/benrowe/nab-bank-api/internal/model"
)

// MockNABClient is a mock implementation of NABClient for testing
type MockNABClient struct{}

// NewMockNABClient creates a new mock NAB client
func NewMockNABClient() NABClient {
	return &MockNABClient{}
}

// GetAccounts returns mock account data
func (m *MockNABClient) GetAccounts(ctx context.Context) ([]model.Account, error) {
	mockAccounts := []model.Account{
		{
			ID:   "12345678",
			Name: "Complete Access Account",
			Type: model.AccountTypeSavings,
			Balance: model.Money{
				Amount: "2543.67",
			},
			AvailableBalance: &model.Money{
				Amount: "2543.67",
			},
			AccountNumber: stringPtr("****5678"),
			BSB:          stringPtr("084001"),
		},
		{
			ID:   "87654321",
			Name: "NAB Classic Banking Account",
			Type: model.AccountTypeChecking,
			Balance: model.Money{
				Amount: "847.23",
			},
			AvailableBalance: &model.Money{
				Amount: "847.23",
			},
			AccountNumber: stringPtr("****4321"),
			BSB:          stringPtr("084001"),
		},
		{
			ID:   "11223344",
			Name: "NAB Reward Saver",
			Type: model.AccountTypeSavings,
			Balance: model.Money{
				Amount: "15420.89",
			},
			AvailableBalance: &model.Money{
				Amount: "15420.89",
			},
			AccountNumber: stringPtr("****3344"),
			BSB:          stringPtr("084001"),
		},
	}

	return mockAccounts, nil
}

// GetAccountTransactions returns mock transaction data
func (m *MockNABClient) GetAccountTransactions(ctx context.Context, accountID string) ([]model.Transaction, error) {
	// Generate some mock transactions based on account ID
	mockTransactions := []model.Transaction{
		{
			ID:          "txn_001_" + accountID,
			Date:        time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
			Description: "EFTPOS Purchase - COLES SUPERMARKET",
			Amount: model.Money{
				Amount: "-85.67",
			},
			Balance: model.Money{
				Amount: "2543.67",
			},
			Category: stringPtr("Groceries"),
			Merchant: stringPtr("COLES SUPERMARKET"),
		},
		{
			ID:          "txn_002_" + accountID,
			Date:        time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
			Description: "Direct Credit - SALARY PAYMENT",
			Amount: model.Money{
				Amount: "2500.00",
			},
			Balance: model.Money{
				Amount: "2629.34",
			},
			Category: stringPtr("Income"),
			Merchant: stringPtr("EMPLOYER PTY LTD"),
		},
		{
			ID:          "txn_003_" + accountID,
			Date:        time.Now().AddDate(0, 0, -3).Format("2006-01-02"),
			Description: "ATM Withdrawal - NAB ATM",
			Amount: model.Money{
				Amount: "-100.00",
			},
			Balance: model.Money{
				Amount: "129.34",
			},
			Category: stringPtr("Cash"),
			Merchant: stringPtr("NAB ATM"),
		},
	}

	return mockTransactions, nil
}

// stringPtr is a helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}