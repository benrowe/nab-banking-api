package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/benrowe/nab-bank-api/internal/model"
	"github.com/benrowe/nab-bank-api/internal/service"
	"github.com/gorilla/mux"
)

// AccountsHandler handles account-related HTTP requests
type AccountsHandler struct {
	accountService service.AccountService
	logger         *log.Logger
}

// NewAccountsHandler creates a new accounts handler
func NewAccountsHandler(accountService service.AccountService, logger *log.Logger) *AccountsHandler {
	return &AccountsHandler{
		accountService: accountService,
		logger:         logger,
	}
}

// ListAccounts handles GET /api/v1/accounts
func (h *AccountsHandler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	h.logger.Printf("ListAccounts: %s %s", r.Method, r.URL.Path)

	accounts, err := h.accountService.GetAllAccounts(r.Context())
	if err != nil {
		h.logger.Printf("Failed to get accounts: %v", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, model.ErrorTypeInternalError, "Failed to retrieve accounts", err)
		return
	}

	response := model.AccountsResponse{
		Accounts:    accounts,
		RetrievedAt: time.Now(),
		Count:       len(accounts),
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// GetAccount handles GET /api/v1/accounts/{accountId}
func (h *AccountsHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["accountId"]
	
	h.logger.Printf("GetAccount: %s %s (ID: %s)", r.Method, r.URL.Path, accountID)

	if accountID == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, model.ErrorTypeInvalidRequest, "Account ID is required", nil)
		return
	}

	accountDetails, err := h.accountService.GetAccountDetails(r.Context(), accountID)
	if err != nil {
		switch err {
		case service.ErrAccountNotFound:
			h.writeErrorResponse(w, http.StatusNotFound, model.ErrorTypeAccountNotFound, "Account not found", nil)
		case service.ErrServiceUnavailable:
			h.writeErrorResponse(w, http.StatusServiceUnavailable, model.ErrorTypeServiceUnavailable, "Service temporarily unavailable", err)
		case service.ErrAuthenticationFailed:
			h.writeErrorResponse(w, http.StatusUnauthorized, model.ErrorTypeAuthenticationFailed, "Authentication failed", nil)
		default:
			h.logger.Printf("Failed to get account details: %v", err)
			h.writeErrorResponse(w, http.StatusInternalServerError, model.ErrorTypeInternalError, "Failed to retrieve account details", err)
		}
		return
	}

	response := model.AccountDetailsResponse{
		Account: *accountDetails,
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// writeJSONResponse writes a JSON response
func (h *AccountsHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Printf("Failed to encode JSON response: %v", err)
	}
}

// writeErrorResponse writes an error response
func (h *AccountsHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, errorType, message string, details interface{}) {
	errorResponse := model.ErrorResponse{
		Error:     errorType,
		Message:   message,
		Details:   details,
		Timestamp: time.Now(),
	}

	h.writeJSONResponse(w, statusCode, errorResponse)
}