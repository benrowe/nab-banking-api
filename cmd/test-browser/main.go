package main

import (
	"context"
	"log"
	"os"

	"github.com/benrowe/nab-bank-api/internal/browser"
	"github.com/benrowe/nab-bank-api/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Override settings for debugging
	cfg.NAB.BrowserHeadless = false // Show browser for debugging
	cfg.NAB.BrowserTimeout = 60     // Longer timeout for manual inspection

	// Create logger
	logger := log.New(os.Stdout, "[NAB-TEST] ", log.LstdFlags|log.Lshortfile)

	// Create NAB client
	nabClient := browser.NewNABClient(&cfg.NAB, logger)

	// Test account retrieval
	ctx := context.Background()
	logger.Println("Starting NAB browser test...")

	accounts, err := nabClient.GetAccounts(ctx)
	if err != nil {
		logger.Printf("Error retrieving accounts: %v", err)
		return
	}

	logger.Printf("Successfully retrieved %d accounts:", len(accounts))
	for i, account := range accounts {
		logger.Printf("  %d. %s (%s) - Balance: $%s", i+1, account.Name, account.ID, account.Balance.Amount)
	}
}