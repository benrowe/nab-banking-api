.PHONY: help

# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

TARGET_MAX_CHAR_NUM=20

# Variables
APP_NAME := nab-bank-api
GO_VERSION := 1.21-alpine
GOLANGCI_LINT_VERSION := latest
DOCKER_REGISTRY ?= 
IMAGE_TAG ?= latest

## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\\-\\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

## Build the Docker image
build:
	@echo "${YELLOW}Building Docker image...${RESET}"
	docker build -t $(APP_NAME):$(IMAGE_TAG) .

## Run the application with Docker
run:
	@echo "${YELLOW}Running application...${RESET}"
	docker run --rm -p 8080:8080 --env-file .env $(APP_NAME):$(IMAGE_TAG)

## Run in development mode with live reload
dev:
	@echo "${YELLOW}Starting development server...${RESET}"
	docker run --rm -p 8080:8080 --env-file .env \
		-v $$(pwd):/app -w /app \
		golang:$(GO_VERSION) go run cmd/server/main.go

## Run all tests
test:
	@echo "${YELLOW}Running tests...${RESET}"
	docker run --rm -v $$(pwd):/app -w /app golang:$(GO_VERSION) go test ./...

## Run tests with coverage
test-coverage:
	@echo "${YELLOW}Running tests with coverage...${RESET}"
	docker run --rm -v $$(pwd):/app -w /app golang:$(GO_VERSION) \
		go test -coverprofile=coverage.out ./...
	docker run --rm -v $$(pwd):/app -w /app golang:$(GO_VERSION) \
		go tool cover -html=coverage.out -o coverage.html

## Run a single test (make test-single TEST=TestName PKG=./pkg)
test-single:
	@echo "${YELLOW}Running single test: $(TEST) in $(PKG)${RESET}"
	docker run --rm -v $$(pwd):/app -w /app golang:$(GO_VERSION) \
		go test -run $(TEST) $(PKG)

## Run golangci-lint
lint:
	@echo "${YELLOW}Running linter...${RESET}"
	docker run --rm -v $$(pwd):/app -w /app \
		golangci/golangci-lint:$(GOLANGCI_LINT_VERSION) \
		golangci-lint run --timeout=5m

## Format Go code
fmt:
	@echo "${YELLOW}Formatting code...${RESET}"
	docker run --rm -v $$(pwd):/app -w /app golang:$(GO_VERSION) go fmt ./...

## Tidy Go modules
tidy:
	@echo "${YELLOW}Tidying modules...${RESET}"
	docker run --rm -v $$(pwd):/app -w /app golang:$(GO_VERSION) go mod tidy

## Initialize Go module (make init MODULE=github.com/user/repo)
init:
	@echo "${YELLOW}Initializing Go module...${RESET}"
	docker run --rm -v $$(pwd):/app -w /app golang:$(GO_VERSION) go mod init $(MODULE)

## Run with delve debugger on port 40000
debug:
	@echo "${YELLOW}Starting debugger on port 40000...${RESET}"
	docker run --rm -v $$(pwd):/app -w /app -p 40000:40000 golang:$(GO_VERSION) \
		sh -c "go install github.com/go-delve/delve/cmd/dlv@latest && \
		dlv debug --headless --listen=:40000 --api-version=2 cmd/server/main.go"

## Open a shell in the Go container
shell:
	@echo "${YELLOW}Opening shell in Go container...${RESET}"
	docker run --rm -it -v $$(pwd):/app -w /app golang:$(GO_VERSION) sh

## Download Go dependencies
deps:
	@echo "${YELLOW}Downloading dependencies...${RESET}"
	docker run --rm -v $$(pwd):/app -w /app golang:$(GO_VERSION) go mod download

## Clean up Docker images and containers
clean:
	@echo "${YELLOW}Cleaning up...${RESET}"
	docker system prune -f
	docker rmi $(APP_NAME):$(IMAGE_TAG) 2>/dev/null || true

## Build and push Docker image to registry
push: build
	@if [ -z "$(DOCKER_REGISTRY)" ]; then \
		echo "${RED}DOCKER_REGISTRY not set${RESET}"; \
		exit 1; \
	fi
	@echo "${YELLOW}Pushing to $(DOCKER_REGISTRY)/$(APP_NAME):$(IMAGE_TAG)${RESET}"
	docker tag $(APP_NAME):$(IMAGE_TAG) $(DOCKER_REGISTRY)/$(APP_NAME):$(IMAGE_TAG)
	docker push $(DOCKER_REGISTRY)/$(APP_NAME):$(IMAGE_TAG)

## Initial project setup
setup:
	@echo "${YELLOW}Setting up project...${RESET}"
	cp .env.example .env
	@echo "${GREEN}Created .env file. Please edit with your credentials.${RESET}"
	@echo "${GREEN}Run 'make init MODULE=your-module-name' to initialize Go module.${RESET}"

## Test browser automation setup
browser-test:
	@echo "${YELLOW}Testing browser automation...${RESET}"
	docker run --rm -v $$(pwd):/app -w /app \
		--cap-add=SYS_ADMIN \
		golang:$(GO_VERSION) \
		sh -c "apk add --no-cache chromium && go run cmd/test-browser/main.go"

## CI-specific linting (exit on failure)
ci-lint:
	docker run --rm -v $$(pwd):/app -w /app \
		golangci/golangci-lint:$(GOLANGCI_LINT_VERSION) \
		golangci-lint run --timeout=5m --issues-exit-code=1

## CI-specific testing with coverage and race detection
ci-test:
	docker run --rm -v $$(pwd):/app -w /app golang:$(GO_VERSION) \
		go test -race -coverprofile=coverage.out -covermode=atomic ./...

## Run security scanning with gosec
security-scan:
	@echo "${YELLOW}Running security scan...${RESET}"
	docker run --rm -v $$(pwd):/app -w /app securecodewarrior/gosec:latest \
		gosec ./...