# Production-ready Multi-stage Docker build for ASA AgriJobs Backend

# Stage 1: Builder
FROM golang:1.24.4-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION:-1.0.0}" \
    -o /app/server ./cmd/server

# Stage 2: Production Runtime
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata curl && \
    rm -rf /var/cache/apk/*

# Create non-root user for security
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder --chown=appuser:appuser /app/server .

# Copy migrations for database setup
COPY --from=builder --chown=appuser:appuser /app/migrations ./migrations

# Copy docs for Swagger API documentation
COPY --from=builder --chown=appuser:appuser /app/docs ./docs

# Switch to non-root user
USER appuser

# Expose application port (configurable via build arg or defaults to 8080)
ARG SERVER_PORT=8080
ENV SERVER_PORT=${SERVER_PORT}
EXPOSE ${SERVER_PORT}

# Health check for container orchestration
# Uses ASA_BASE_URL from environment for health monitoring
HEALTHCHECK --interval=30s --timeout=5s --start-period=60s --retries=3 \
  CMD sh -c 'curl -f ${ASA_BASE_URL}/health || exit 1'

# Run the application
ENTRYPOINT ["./server"]
