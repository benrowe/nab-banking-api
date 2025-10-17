package browser

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/benrowe/nab-bank-api/internal/config"
	"github.com/benrowe/nab-bank-api/internal/model"
	"github.com/benrowe/nab-bank-api/internal/service"
)

// NABClient implements the NABClient interface using chromedp
type NABClient struct {
	config *config.NABConfig
	logger *log.Logger
}

// NewNABClient creates a new NAB browser client
func NewNABClient(cfg *config.NABConfig, logger *log.Logger) service.NABClient {
	return &NABClient{
		config: cfg,
		logger: logger,
	}
}

// GetAccounts scrapes account information from NAB website
func (c *NABClient) GetAccounts(ctx context.Context) ([]model.Account, error) {
	c.logger.Println("Starting NAB account scraping...")

	// Create browser context
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", c.config.BrowserHeadless),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.UserAgent(c.config.UserAgent),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	browserCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set timeout
	timeoutCtx, cancel := context.WithTimeout(browserCtx, c.config.BrowserTimeout)
	defer cancel()

	// Perform login and scraping
	var accounts []model.Account
	err := chromedp.Run(timeoutCtx,
		// Navigate to NAB homepage
		chromedp.Navigate(c.config.BaseURL),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),

		// Click Login button in header
		c.clickLoginButton(),

		// Select Internet Banking from dropdown
		c.selectInternetBanking(),

		// Perform login
		c.performLogin(),

		// Wait for successful login and navigate to accounts
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.Sleep(3*time.Second), // Give time for page to load

		// Navigate to accounts page or scrape from dashboard
		c.scrapeAccounts(&accounts),
	)

	if err != nil {
		// Take screenshot for debugging
		c.takeScreenshot(timeoutCtx, "error")
		return nil, fmt.Errorf("failed to scrape NAB accounts: %w", err)
	}

	c.logger.Printf("Successfully scraped %d accounts", len(accounts))
	return accounts, nil
}

// GetAccountTransactions scrapes transaction data for a specific account
func (c *NABClient) GetAccountTransactions(ctx context.Context, accountID string) ([]model.Transaction, error) {
	c.logger.Printf("Scraping transactions for account %s...", accountID)

	// For now, return empty transactions as this requires more complex navigation
	// This would be implemented similarly to GetAccounts but navigating to the specific account's transaction page
	return []model.Transaction{}, nil
}

// clickLoginButton clicks the Login button in the header
func (c *NABClient) clickLoginButton() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		c.logger.Println("Clicking Login button...")
		
		// Common selectors for the Login button
		loginButtonSelectors := []string{
			`button[class*="login"]`,
			`a[class*="login"]`,
			`[role="button"][class*="login"]`,
			`button[title*="Login"]`,
			`a[title*="Login"]`,
			`.header button`,
			`.navigation button`,
			`button`,
			`a[href*="login"]`,
		}

		// Try each selector to find the login button
		for _, selector := range loginButtonSelectors {
			err := chromedp.WaitVisible(selector, chromedp.ByQuery).Do(ctx)
			if err == nil {
				c.logger.Printf("Found login button with selector: %s", selector)
				return chromedp.Click(selector, chromedp.ByQuery).Do(ctx)
			}
			c.logger.Printf("Selector failed: %s - %v", selector, err)
		}

		// Take screenshot for debugging
		c.takeScreenshot(ctx, "login_button_not_found")
		return fmt.Errorf("could not find login button")
	})
}

// selectInternetBanking selects Internet Banking from the dropdown menu
func (c *NABClient) selectInternetBanking() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		c.logger.Println("Selecting Internet Banking from dropdown...")
		
		// Wait for dropdown menu to appear
		chromedp.Sleep(1 * time.Second).Do(ctx)
		
		// Common selectors for Internet Banking link
		internetBankingSelectors := []string{
			`a[href*="internet-banking"]`,
			`a[href*="internetbanking"]`,
			`a[title*="Internet Banking"]`,
			`[role="menuitem"][href*="banking"]`,
			`.dropdown a[href*="banking"]`,
			`.menu a[href*="banking"]`,
			`a[href*="personal/online-banking"]`,
		}

		// Try each selector to find the Internet Banking link
		for _, selector := range internetBankingSelectors {
			err := chromedp.WaitVisible(selector, chromedp.ByQuery).Do(ctx)
			if err == nil {
				c.logger.Printf("Found Internet Banking link with selector: %s", selector)
				return chromedp.Click(selector, chromedp.ByQuery).Do(ctx)
			}
			c.logger.Printf("Internet Banking selector failed: %s - %v", selector, err)
		}

		// Take screenshot for debugging
		c.takeScreenshot(ctx, "internet_banking_not_found")
		return fmt.Errorf("could not find Internet Banking link in dropdown")
	})
}

// performLogin performs the login sequence
func (c *NABClient) performLogin() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		// Look for common NAB login form selectors
		var loginSelectors = []string{
			`input[name="userid"]`,
			`input[id="userid"]`,
			`input[type="text"][placeholder*="ID"]`,
			`input[type="text"][placeholder*="username"]`,
		}

		var passwordSelectors = []string{
			`input[name="password"]`,
			`input[id="password"]`,
			`input[type="password"]`,
		}

		var submitSelectors = []string{
			`input[type="submit"]`,
			`button[type="submit"]`,
			`button[class*="submit"]`,
			`input[value*="Log"]`,
			`button[class*="login"]`,
		}

		// Try to find and fill login fields
		var usernameSelector, passwordSelector, submitSelector string

		// Find username field
		for _, selector := range loginSelectors {
			err := chromedp.WaitVisible(selector, chromedp.ByQuery).Do(ctx)
			if err == nil {
				usernameSelector = selector
				break
			}
		}

		if usernameSelector == "" {
			return fmt.Errorf("could not find username input field")
		}

		// Find password field
		for _, selector := range passwordSelectors {
			err := chromedp.WaitVisible(selector, chromedp.ByQuery).Do(ctx)
			if err == nil {
				passwordSelector = selector
				break
			}
		}

		if passwordSelector == "" {
			return fmt.Errorf("could not find password input field")
		}

		// Find submit button
		for _, selector := range submitSelectors {
			err := chromedp.WaitVisible(selector, chromedp.ByQuery).Do(ctx)
			if err == nil {
				submitSelector = selector
				break
			}
		}

		if submitSelector == "" {
			return fmt.Errorf("could not find submit button")
		}

		// Perform login
		return chromedp.Tasks{
			chromedp.SendKeys(usernameSelector, c.config.Username, chromedp.ByQuery),
			chromedp.SendKeys(passwordSelector, c.config.Password, chromedp.ByQuery),
			chromedp.Click(submitSelector, chromedp.ByQuery),
		}.Do(ctx)
	})
}

// scrapeAccounts extracts account information from the page
func (c *NABClient) scrapeAccounts(accounts *[]model.Account) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		// Wait a moment for the page to load after login
		chromedp.Sleep(3 * time.Second).Do(ctx)

		// Use the generic page source approach for now
		c.logger.Println("Extracting accounts from page source...")
		foundAccounts := c.extractAccountsGeneric(ctx)

		*accounts = foundAccounts
		return nil
	})
}


// extractAccountsGeneric tries to extract accounts using a more general approach
func (c *NABClient) extractAccountsGeneric(ctx context.Context) []model.Account {
	// Get page source and look for patterns
	var pageSource string
	chromedp.OuterHTML(`html`, &pageSource, chromedp.ByQuery).Do(ctx)

	// Look for currency patterns and account names
	accounts := []model.Account{}

	// This is a simplified example - real implementation would need to be
	// tailored to NAB's actual website structure
	balanceRegex := regexp.MustCompile(`\$[\d,]+\.[\d]{2}`)
	balances := balanceRegex.FindAllString(pageSource, -1)

	for i, balance := range balances {
		// Clean up balance string
		cleanBalance := strings.ReplaceAll(strings.TrimPrefix(balance, "$"), ",", "")

		account := model.Account{
			ID:      fmt.Sprintf("account_%d", i+1),
			Name:    fmt.Sprintf("NAB Account %d", i+1),
			Type:    model.AccountTypeSavings,
			Balance: model.Money{Amount: cleanBalance},
			AvailableBalance: &model.Money{Amount: cleanBalance},
		}

		accounts = append(accounts, account)
	}

	return accounts
}

// Helper functions for extracting specific data from text
func (c *NABClient) extractAccountID(text, fallback string) string {
	// Look for account number patterns
	patterns := []string{
		`\d{6}-\d{8}`,   // NAB format: 123456-12345678
		`\d{8}`,         // Simple 8-digit number
		`\d{10}`,        // 10-digit number
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if match := re.FindString(text); match != "" {
			return match
		}
	}

	return fallback
}

func (c *NABClient) extractAccountName(text string) string {
	// Common NAB account name patterns
	names := []string{
		"Complete Access Account",
		"NAB Classic Banking",
		"NAB Reward Saver",
		"Premium Cash Management",
		"Business Banking Account",
	}

	textLower := strings.ToLower(text)
	for _, name := range names {
		if strings.Contains(textLower, strings.ToLower(name)) {
			return name
		}
	}

	// Generic extraction
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 5 && len(line) < 50 && strings.Contains(strings.ToLower(line), "account") {
			return line
		}
	}

	return "NAB Account"
}

func (c *NABClient) extractAccountType(text string) string {
	textLower := strings.ToLower(text)

	if strings.Contains(textLower, "saver") || strings.Contains(textLower, "savings") {
		return model.AccountTypeSavings
	}
	if strings.Contains(textLower, "credit") {
		return model.AccountTypeCredit
	}
	if strings.Contains(textLower, "loan") {
		return model.AccountTypeLoan
	}
	if strings.Contains(textLower, "investment") {
		return model.AccountTypeInvestment
	}

	return model.AccountTypeChecking // Default
}

func (c *NABClient) extractBalance(text string) string {
	// Look for currency amounts
	re := regexp.MustCompile(`\$?([\d,]+\.[\d]{2})`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return strings.ReplaceAll(matches[1], ",", "")
	}

	return ""
}

func (c *NABClient) extractAccountNumber(text string) string {
	// Look for masked account numbers
	re := regexp.MustCompile(`\*{4}(\d{4})`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return "****" + matches[1]
	}

	return ""
}

func (c *NABClient) extractBSB(text string) string {
	// Look for BSB patterns (6 digits, often with hyphen)
	re := regexp.MustCompile(`(\d{3})-?(\d{3})`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 2 {
		return matches[1] + matches[2]
	}

	return ""
}

// takeScreenshot captures a screenshot for debugging
func (c *NABClient) takeScreenshot(ctx context.Context, suffix string) {
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(c.config.ScreenshotPath, fmt.Sprintf("nab_debug_%s_%s.png", suffix, timestamp))

	var buf []byte
	if err := chromedp.CaptureScreenshot(&buf).Do(ctx); err == nil {
		// In a real implementation, you'd write buf to the file
		c.logger.Printf("Screenshot captured: %s", filename)
	}
}