# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Run stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create necessary directories
RUN mkdir -p settings logs storage assets

# Copy binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 5000

# Run the application
CMD ["./main"]
