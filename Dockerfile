# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Create non-root user
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /build

# Copy go mod files first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o nab-bank-api \
    cmd/server/main.go

# Development stage
FROM golang:1.21-alpine AS development

# Install development tools and browser dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    chromium \
    nss \
    freetype \
    freetype-dev \
    harfbuzz \
    ttf-freefont \
    make \
    curl

# Install air for live reload
RUN go install github.com/cosmtrek/air@v1.49.0

# Install delve debugger
RUN go install github.com/go-delve/delve/cmd/dlv@v1.21.2

# Install golangci-lint
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2

# Create non-root user
RUN adduser -D -g '' appuser

# Set Chrome executable path
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROME_PATH=/usr/bin/chromium-browser

# Create directories for screenshots and downloads
RUN mkdir -p /app/screenshots /app/downloads && \
    chown -R appuser:appuser /app

# Set working directory
WORKDIR /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080
EXPOSE 40000

# Default command for development
CMD ["air", "-c", ".air.toml"]

# Production stage
FROM alpine:3.18 AS production

# Install runtime dependencies for browser automation
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ttf-freefont \
    dumb-init

# Create non-root user
RUN adduser -D -g '' appuser

# Set Chrome executable path
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROME_PATH=/usr/bin/chromium-browser

# Create directories for screenshots and downloads
RUN mkdir -p /app/screenshots /app/downloads && \
    chown -R appuser:appuser /app

# Copy the binary from builder stage
COPY --from=builder /build/nab-bank-api /app/nab-bank-api
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Set permissions
RUN chown appuser:appuser /app/nab-bank-api && \
    chmod +x /app/nab-bank-api

# Switch to non-root user
USER appuser

# Set working directory
WORKDIR /app

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Use dumb-init as PID 1 to handle signals properly
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

# Run the binary
CMD ["./nab-bank-api"]