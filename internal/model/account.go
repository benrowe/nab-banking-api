package model

import (
	"time"
)

// Money represents a monetary amount
type Money struct {
	Amount string `json:"amount" example:"1234.56"`
}

// Account represents a bank account
type Account struct {
	ID               string     `json:"id" example:"12345678"`
	Name             string     `json:"name" example:"Complete Access Account"`
	Type             string     `json:"type" example:"savings"`
	Balance          Money      `json:"balance"`
	AvailableBalance *Money     `json:"availableBalance,omitempty"`
	AccountNumber    *string    `json:"accountNumber,omitempty" example:"****1234"`
	BSB              *string    `json:"bsb,omitempty" example:"084001"`
	LastUpdated      *time.Time `json:"lastUpdated,omitempty"`
}

// AccountsResponse represents the response for listing accounts
type AccountsResponse struct {
	Accounts    []Account `json:"accounts"`
	RetrievedAt time.Time `json:"retrievedAt"`
	Count       int       `json:"count" example:"3"`
}

// Transaction represents a bank transaction
type Transaction struct {
	ID          string     `json:"id" example:"txn_20231017_001"`
	Date        string     `json:"date" example:"2023-10-17"`
	Description string     `json:"description" example:"EFTPOS Purchase - COLES SUPERMARKET"`
	Amount      Money      `json:"amount"`
	Balance     Money      `json:"balance"`
	Category    *string    `json:"category,omitempty" example:"Groceries"`
	Merchant    *string    `json:"merchant,omitempty" example:"COLES SUPERMARKET"`
}

// AccountDetails extends Account with transaction information
type AccountDetails struct {
	Account
	Transactions             []Transaction `json:"transactions,omitempty"`
	RecentTransactionCount   int           `json:"recentTransactionCount,omitempty" example:"10"`
}

// AccountDetailsResponse represents the response for getting account details
type AccountDetailsResponse struct {
	Account AccountDetails `json:"account"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error     string      `json:"error" example:"AUTHENTICATION_FAILED"`
	Message   string      `json:"message" example:"Invalid credentials provided"`
	Details   interface{} `json:"details,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// AccountType constants
const (
	AccountTypeSavings    = "savings"
	AccountTypeChecking   = "checking"
	AccountTypeCredit     = "credit"
	AccountTypeLoan       = "loan"
	AccountTypeInvestment = "investment"
)

// Error types
const (
	ErrorTypeAuthenticationFailed = "AUTHENTICATION_FAILED"
	ErrorTypeAccountNotFound      = "ACCOUNT_NOT_FOUND"
	ErrorTypeServiceUnavailable   = "SERVICE_UNAVAILABLE"
	ErrorTypeInternalError        = "INTERNAL_ERROR"
	ErrorTypeInvalidRequest       = "INVALID_REQUEST"
)