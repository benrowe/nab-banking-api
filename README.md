# NAB Bank Automation

A Go-based microservice for automating NAB Bank interactions using headless browser automation, built with Docker and deployed via GitHub Actions.

## Quick Start

### Prerequisites
- Docker
- Docker Compose (optional)
- Git

### Development Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd nab-bank-api
```

2. Set up the project:
```bash
make setup
# Edit .env with your NAB banking credentials
```

3. Initialize Go module:
```bash
make init MODULE=github.com/yourusername/nab-bank-api
```

4. Build and run the application:
```bash
make build
make run
```

### Development Commands

```bash
# Show all available commands
make help

# Development workflow
make dev          # Run in development mode with live reload
make test         # Run all tests
make lint         # Run linting
make fmt          # Format code

# Build and deployment
make build        # Build Docker image
make run          # Run the application
```

See [WARP.md](WARP.md) for comprehensive development commands and guidelines.

## Architecture

This service follows clean architecture principles with the following structure:

- `cmd/server/` - Application entry point
- `internal/api/` - HTTP handlers and routing
- `internal/service/` - Business logic
- `internal/browser/` - Browser automation client
- `internal/pages/` - Page object models for NAB web interface
- `internal/model/` - Data models
- `internal/config/` - Configuration management

## API Endpoints

- `GET /health` - Health check endpoint
- `GET /ready` - Readiness check endpoint

## Configuration

Environment variables:
- `NAB_USERNAME` - NAB banking username
- `NAB_PASSWORD` - NAB banking password
- `NAB_BASE_URL` - NAB website URL (default: https://www.nab.com.au)
- `BROWSER_HEADLESS` - Run browser in headless mode (default: true)
- `BROWSER_TIMEOUT` - Browser operation timeout in seconds (default: 30)
- `PORT` - Server port (default: 8080)
- `LOG_LEVEL` - Log level (default: info)

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run a specific test
make test-single TEST=TestFunctionName PKG=./internal/service

# Test browser automation setup
make browser-test
```

## Deployment

This project uses GitHub Actions for CI/CD. See `.github/workflows/` for pipeline configuration.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

[Add your license here]