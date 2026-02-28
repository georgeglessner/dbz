# Build stage
FROM golang:1.21-alpine AS builder

# Accept version as build argument
ARG VERSION=dev

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary with version
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s -X main.Version=${VERSION}" -o dbz main.go

# Final stage
FROM alpine:latest

# Accept version for label
ARG VERSION=dev

# Install runtime dependencies
RUN apk add --no-cache ca-certificates docker-cli

# Create non-root user
RUN addgroup -g 1000 dbz && \
    adduser -D -u 1000 -G dbz dbz

# Set working directory
WORKDIR /home/dbz

# Copy binary from builder
COPY --from=builder /app/dbz /usr/local/bin/dbz

# Make binary executable
RUN chmod +x /usr/local/bin/dbz

# Switch to non-root user
USER dbz

# Set entrypoint
ENTRYPOINT ["dbz"]

# Default command
CMD ["--help"]

# Labels
LABEL maintainer="George Glessner"
LABEL description="dbz - Database CLI Tool"
LABEL version="${VERSION}"

# Expose common database ports (for reference)
EXPOSE 5432 3306 3307 8123