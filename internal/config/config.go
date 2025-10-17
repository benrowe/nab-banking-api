package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	Server ServerConfig
	NAB    NABConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string
}

// NABConfig holds NAB-specific configuration
type NABConfig struct {
	Username        string
	Password        string
	BaseURL         string
	LoginURL        string
	AccountsURL     string
	BrowserTimeout  time.Duration
	BrowserHeadless bool
	ScreenshotPath  string
	UserAgent       string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port: getEnvOrDefault("PORT", "8080"),
		},
		NAB: NABConfig{
			Username:        os.Getenv("NAB_USERNAME"),
			Password:        os.Getenv("NAB_PASSWORD"),
			BaseURL:         getEnvOrDefault("NAB_BASE_URL", "https://www.nab.com.au"),
			LoginURL:        getEnvOrDefault("NAB_LOGIN_URL", "https://www.nab.com.au/personal/online-banking/nab-internet-banking"),
			AccountsURL:     getEnvOrDefault("NAB_ACCOUNTS_URL", "/internetbanking/AccountBalance.jsp"),
			BrowserTimeout:  parseDurationOrDefault("BROWSER_TIMEOUT", 30*time.Second),
			BrowserHeadless: parseBoolOrDefault("BROWSER_HEADLESS", true),
			ScreenshotPath:  getEnvOrDefault("BROWSER_SCREENSHOT_PATH", "/app/screenshots"),
			UserAgent:       getEnvOrDefault("BROWSER_USER_AGENT", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
		},
	}

	// Validate required fields
	if config.NAB.Username == "" {
		return nil, fmt.Errorf("NAB_USERNAME environment variable is required")
	}
	if config.NAB.Password == "" {
		return nil, fmt.Errorf("NAB_PASSWORD environment variable is required")
	}

	return config, nil
}

// getEnvOrDefault gets environment variable value or returns default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseDurationOrDefault parses duration from env var or returns default
func parseDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// parseBoolOrDefault parses boolean from env var or returns default
func parseBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolean, err := strconv.ParseBool(value); err == nil {
			return boolean
		}
	}
	return defaultValue
}