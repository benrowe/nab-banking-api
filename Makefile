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

## Build the Docker image for production
build:
	@echo "${YELLOW}Building production Docker image...${RESET}"
	docker build --target production -t $(APP_NAME):$(IMAGE_TAG) .

## Build the Docker image for development
build-dev:
	@echo "${YELLOW}Building development Docker image...${RESET}"
	docker build --target development -t $(APP_NAME):dev .

## Run the application with Docker (production)
run:
	@echo "${YELLOW}Running application...${RESET}"
	docker run --rm -p 8080:8080 --env-file .env $(APP_NAME):$(IMAGE_TAG)

## Run in development mode with live reload
dev: build-dev
	@echo "${YELLOW}Starting development server with live reload...${RESET}"
	docker run --rm -it -p 8080:8080 -p 40000:40000 --env-file .env \
		-v $$(pwd):/app \
		$(APP_NAME):dev

## Run all tests
test: build-dev
	@echo "${YELLOW}Running tests...${RESET}"
	docker run --rm -v $$(pwd):/app --user root $(APP_NAME):dev go test ./...

## Run tests with coverage
test-coverage: build-dev
	@echo "${YELLOW}Running tests with coverage...${RESET}"
	docker run --rm -v $$(pwd):/app $(APP_NAME):dev \
		go test -coverprofile=coverage.out ./...
	docker run --rm -v $$(pwd):/app $(APP_NAME):dev \
		go tool cover -html=coverage.out -o coverage.html

## Run a single test (make test-single TEST=TestName PKG=./pkg)
test-single: build-dev
	@echo "${YELLOW}Running single test: $(TEST) in $(PKG)${RESET}"
	docker run --rm -v $$(pwd):/app $(APP_NAME):dev \
		go test -run $(TEST) $(PKG)

## Run golangci-lint
lint: build-dev
	@echo "${YELLOW}Running linter...${RESET}"
	docker run --rm -v $$(pwd):/app $(APP_NAME):dev \
		golangci-lint run --timeout=5m

## Format Go code
fmt: build-dev
	@echo "${YELLOW}Formatting code...${RESET}"
	docker run --rm -v $$(pwd):/app $(APP_NAME):dev go fmt ./...

## Tidy Go modules
tidy: build-dev
	@echo "${YELLOW}Tidying modules...${RESET}"
	docker run --rm -v $$(pwd):/app $(APP_NAME):dev go mod tidy

## Initialize Go module (make init MODULE=github.com/user/repo)
init: build-dev
	@echo "${YELLOW}Initializing Go module...${RESET}"
	docker run --rm -v $$(pwd):/app $(APP_NAME):dev go mod init $(MODULE)

## Run with delve debugger on port 40000
debug: build-dev
	@echo "${YELLOW}Starting debugger on port 40000...${RESET}"
	docker run --rm -it -v $$(pwd):/app -p 40000:40000 --env-file .env \
		$(APP_NAME):dev \
		dlv debug --headless --listen=:40000 --api-version=2 cmd/server/main.go

## Open a shell in the Go container
shell: build-dev
	@echo "${YELLOW}Opening shell in Go container...${RESET}"
	docker run --rm -it -v $$(pwd):/app $(APP_NAME):dev sh

## Download Go dependencies
deps: build-dev
	@echo "${YELLOW}Downloading dependencies...${RESET}"
	docker run --rm -v $$(pwd):/app $(APP_NAME):dev go mod download

## Clean up Docker images and containers
clean:
	@echo "${YELLOW}Cleaning up...${RESET}"
	docker system prune -f
	docker rmi $(APP_NAME):$(IMAGE_TAG) 2>/dev/null || true
	docker rmi $(APP_NAME):dev 2>/dev/null || true

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
browser-test: build-dev
	@echo "${YELLOW}Testing browser automation...${RESET}"
	docker run --rm -v $$(pwd):/app --cap-add=SYS_ADMIN --env-file .env \
		$(APP_NAME):dev go run cmd/test-browser/main.go

## CI-specific linting (exit on failure)
ci-lint: build-dev
	docker run --rm -v $$(pwd):/app \
		$(APP_NAME):dev \
		golangci-lint run --timeout=5m --issues-exit-code=1

## CI-specific testing with coverage and race detection
ci-test: build-dev
	docker run --rm -v $$(pwd):/app $(APP_NAME):dev \
		go test -race -coverprofile=coverage.out -covermode=atomic ./...

## Run security scanning with gosec
security-scan: build-dev
	@echo "${YELLOW}Running security scan...${RESET}"
	docker run --rm -v $$(pwd):/app $(APP_NAME):dev \
		sh -c "go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest && gosec ./..."
