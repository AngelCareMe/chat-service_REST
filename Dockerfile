# Dockerfile
FROM golang:1.24.4-alpine AS builder

# Install git and build tools
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /app

# Copy everything (вместо отдельных файлов)
COPY . .

# Check if files exist
RUN ls -la go.mod go.sum

# Download dependencies
RUN go mod download

# Build the binary
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM alpine:latest

# Install ca-certificates and wget for health check
RUN apk --no-cache add ca-certificates wget

# Create non-root user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy configs and migrations
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/migrations ./migrations

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# Command to run the application
CMD ["./main"]