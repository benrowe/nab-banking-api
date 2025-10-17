package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/benrowe/nab-bank-api/internal/api/handler"
	"github.com/benrowe/nab-bank-api/internal/browser"
	"github.com/benrowe/nab-bank-api/internal/config"
	"github.com/benrowe/nab-bank-api/internal/service"
	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize dependencies
	logger := log.New(os.Stdout, "[NAB-API] ", log.LstdFlags|log.Lshortfile)
	
	// Choose client based on environment
	var nabClient service.NABClient
	if cfg.NAB.Username == "test" && cfg.NAB.Password == "test" {
		// Use mock client for testing
		logger.Println("Using mock NAB client for testing")
		nabClient = service.NewMockNABClient()
	} else {
		// Use real browser client
		logger.Println("Using real NAB browser client")
		nabClient = browser.NewNABClient(&cfg.NAB, logger)
	}
	
	accountService := service.NewAccountService(nabClient)
	accountsHandler := handler.NewAccountsHandler(accountService, logger)

	// Setup routes
	router := mux.NewRouter()
	
	// Health check
	router.HandleFunc("/health", healthHandler).Methods("GET")
	
	// Hello world (for backward compatibility)
	router.HandleFunc("/", helloHandler).Methods("GET")
	
	// API v1 routes
	v1 := router.PathPrefix("/api/v1").Subrouter()
	v1.HandleFunc("/accounts", accountsHandler.ListAccounts).Methods("GET")
	v1.HandleFunc("/accounts/{accountId}", accountsHandler.GetAccount).Methods("GET")

	// Add middleware
	router.Use(loggingMiddleware(logger))
	router.Use(corsMiddleware)

	logger.Printf("Server starting on port %s", cfg.Server.Port)
	logger.Printf("API endpoints:")
	logger.Printf("  GET /health - Health check")
	logger.Printf("  GET /api/v1/accounts - List all accounts")
	logger.Printf("  GET /api/v1/accounts/{id} - Get account details")
	
	if err := http.ListenAndServe(":"+cfg.Server.Port, router); err != nil {
		log.Fatal(err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World! NAB Bank API is running.\n")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK\n")
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(logger *log.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
			next.ServeHTTP(w, r)
		})
	}
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}
