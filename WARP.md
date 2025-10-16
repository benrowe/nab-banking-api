# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Project Overview

This is a Go-based NAB Bank automation service that uses headless browser automation to interact with NAB's web banking interface. The service runs in Docker containers with GitHub Actions CI/CD pipeline. No local Go installation required - all development happens through Docker.

## Development Commands

### Quick Start
```bash
# Initial project setup
make setup        # Creates .env from .env.example
make init MODULE=github.com/yourusername/nab-bank-api  # Initialize Go module

# Daily development workflow
make help         # Show all available commands
make dev          # Start development server with live reload
make test         # Run all tests
make lint         # Run linting
make fmt          # Format code
```

### Dockerfile-based Development Workflow
```bash
# Build production image
make build
# Or: docker build --target production -t nab-bank-api .

# Build development image
make build-dev
# Or: docker build --target development -t nab-bank-api:dev .

# Run the production application
make run
# Or: docker run -p 8080:8080 --env-file .env nab-bank-api

# Run in development mode with live reload (using Air)
make dev
# Or: docker run --rm -it -p 8080:8080 -p 40000:40000 --env-file .env -v $(pwd):/app nab-bank-api:dev

# Run tests (automatically builds dev image)
make test
# Or: docker run --rm -v $(pwd):/app nab-bank-api:dev go test ./...

# Run tests with coverage
make test-coverage
# Or: docker run --rm -v $(pwd):/app nab-bank-api:dev go test -coverprofile=coverage.out ./...

# Run linting (using built-in golangci-lint)
make lint
# Or: docker run --rm -v $(pwd):/app nab-bank-api:dev golangci-lint run --timeout=5m

# Format code
make fmt
# Or: docker run --rm -v $(pwd):/app nab-bank-api:dev go fmt ./...

# Tidy dependencies
make tidy
# Or: docker run --rm -v $(pwd):/app nab-bank-api:dev go mod tidy

# Run specific test
make test-single TEST=TestFunctionName PKG=./internal/service
# Or: docker run --rm -v $(pwd):/app nab-bank-api:dev go test -run TestFunctionName ./internal/service
```

### Docker Compose (if using)
```bash
# Start all services
docker-compose up

# Build and start
docker-compose up --build

# Run tests
docker-compose run --rm app go test ./...

# Stop services
docker-compose down
```

### Makefile Commands
```bash
# Show all available commands with descriptions
make help

# Development workflow
make dev          # Run in development mode with live reload (auto-builds dev image)
make test         # Run all tests (auto-builds dev image)
make test-coverage # Run tests with coverage report (auto-builds dev image)
make test-single  # Run specific test (TEST=TestName PKG=./pkg)
make lint         # Run golangci-lint (auto-builds dev image)
make fmt          # Format Go code (auto-builds dev image)
make tidy         # Tidy Go modules (auto-builds dev image)

# Build and deployment
make build        # Build production Docker image
make build-dev    # Build development Docker image
make run          # Run the production application
make push         # Build and push to registry

# Project setup
make setup        # Initial project setup (creates .env from example)
make init         # Initialize Go module (MODULE=github.com/user/repo)
make deps         # Download Go dependencies (auto-builds dev image)

# Development tools
make debug        # Run with delve debugger on port 40000 (auto-builds dev image)
make shell        # Open shell in Go container (auto-builds dev image)
make browser-test # Test browser automation setup (auto-builds dev image)

# CI/CD commands
make ci-lint      # CI-specific linting (exit on failure)
make ci-test      # CI testing with race detection and coverage
make security-scan # Run gosec security scanning (auto-builds dev image)

# Cleanup
make clean        # Clean up Docker images and containers (both prod and dev)
```

## Architecture & Project Structure

### Standard Go Project Layout
```
├── cmd/
│   └── server/          # Application entrypoints
├── internal/
│   ├── api/            # API handlers and routes
│   ├── service/        # Business logic
│   ├── browser/        # Browser automation client
│   ├── pages/          # Page object models
│   ├── model/          # Data models
│   ├── config/         # Configuration
│   └── middleware/     # HTTP middleware
├── pkg/                # Public libraries (if any)
├── deployments/        # Docker, k8s configs
├── .github/
│   └── workflows/      # GitHub Actions
├── Dockerfile          # Multi-stage build (development + production)
├── .air.toml          # Live reload configuration
├── docker-compose.yml
└── go.mod
```

### Key Architectural Principles
- **Clean Architecture**: Separate concerns between handlers, services, and data access
- **Dependency Injection**: Use interfaces for testability
- **Configuration**: Environment-based config with validation
- **Error Handling**: Structured error handling with proper HTTP status codes
- **Logging**: Structured logging with correlation IDs
- **Health Checks**: Implement `/health` and `/ready` endpoints

### NAB Browser Automation Patterns
- **Browser Interface**: Abstract browser operations behind interfaces
- **Page Object Model**: Organize page interactions into reusable objects
- **Wait Strategies**: Implement explicit waits for dynamic content
- **Error Handling**: Handle browser timeouts, element not found, and navigation errors
- **Screenshot Capture**: Take screenshots on failures for debugging
- **Session Management**: Handle login sessions and re-authentication
- **Headless Mode**: Run browsers without GUI for production environments

## GitHub Actions CI/CD

### Workflow Triggers
- Push to `main` branch triggers deployment
- Pull requests trigger testing and linting
- Manual workflow dispatch available

### Pipeline Stages
1. **Lint & Format Check**: golangci-lint, gofmt
2. **Test**: Unit tests with coverage reporting
3. **Security Scan**: gosec, nancy for dependency scanning
4. **Build**: Multi-stage Docker build
5. **Deploy**: Deploy to staging/production

### Environment Variables for CI
```yaml
# Required secrets in GitHub
NAB_USERNAME
NAB_PASSWORD
NAB_BASE_URL
DOCKER_REGISTRY_URL
DOCKER_USERNAME
DOCKER_PASSWORD
```

## Development Best Practices

### Go Code Standards
- Use `gofmt` for formatting
- Follow effective Go guidelines
- Use meaningful package names (avoid `util`, `common`)
- Implement proper error handling (don't ignore errors)
- Use context.Context for cancellation and timeouts
- Validate all external inputs

### Docker Best Practices
- **Multi-stage builds**: Dockerfile includes separate builder, development, and production stages
- **Non-root user**: All stages run as non-root user for security
- **Alpine base images**: Minimal attack surface and smaller image size
- **Optimized caching**: Go modules downloaded before copying source code
- **Development tools**: Development stage includes Air for live reload and Delve debugger
- **Browser automation**: All stages include Chromium and required dependencies
- **Health checks**: Production stage includes health check endpoint monitoring

### Testing Strategy
- Unit tests for all business logic
- Integration tests for NAB API client
- Use testify for assertions and mocking
- Mock external dependencies
- Test error scenarios and edge cases
- Maintain >80% code coverage

### Security Considerations
- Store secrets in environment variables only
- Never commit API keys or credentials
- Use HTTPS for all external API calls
- Implement request timeout and context cancellation
- Validate and sanitize all inputs
- Use structured logging without exposing sensitive data

## Environment Configuration

### Required Environment Variables
```bash
# NAB Banking Credentials
NAB_USERNAME=your-nab-username
NAB_PASSWORD=your-nab-password
NAB_BASE_URL=https://www.nab.com.au

# Browser Configuration
BROWSER_HEADLESS=true
BROWSER_TIMEOUT=30
BROWSER_SCREENSHOT_PATH=/app/screenshots
BROWSER_DOWNLOADS_PATH=/app/downloads

# Application Configuration
PORT=8080
LOG_LEVEL=info
ENVIRONMENT=development
```

### Docker Environment File (.env)
Create `.env` file for local development:
```bash
# Use make setup for initial project setup
make setup
# This creates .env from .env.example - edit with your values

# Or manually:
cp .env.example .env
# Edit .env with your values
```

## Common Go Dependencies
```go
// HTTP framework
github.com/gin-gonic/gin
github.com/gorilla/mux

// Browser automation
github.com/chromedp/chromedp
github.com/playwright-community/playwright-go

// Configuration
github.com/spf13/viper
github.com/kelseyhightower/envconfig

// Logging
github.com/sirupsen/logrus
go.uber.org/zap

// Testing
github.com/stretchr/testify
github.com/golang/mock

// HTML parsing
github.com/PuerkitoBio/goquery
golang.org/x/net/html

// Database (if needed)
gorm.io/gorm
github.com/lib/pq
```
### Debugging in Docker

### VS Code with Docker
- Use Remote-Containers extension
- Configure `.devcontainer/devcontainer.json`
- Enable Go debugging in container

### Debug Commands
```bash
# Run with delve debugger (uses development image with pre-installed debugger)
make debug
# Or: docker run --rm -it -v $(pwd):/app -p 40000:40000 --env-file .env nab-bank-api:dev dlv debug --headless --listen=:40000 --api-version=2 cmd/server/main.go

# Shell into development container
make shell
# Or: docker run --rm -it -v $(pwd):/app nab-bank-api:dev sh

# Run development server with live reload
make dev
# This uses Air for automatic recompilation on file changes
```
```