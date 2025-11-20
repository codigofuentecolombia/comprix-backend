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

# Install ca-certificates and bash for HTTPS
RUN apk --no-cache add ca-certificates bash

# Create necessary directories
RUN mkdir -p settings logs storage assets

# Copy binary from builder
COPY --from=builder /app/main .

# Create entrypoint script inline
RUN echo '#!/bin/bash' > /app/entrypoint.sh && \
    echo 'set -e' >> /app/entrypoint.sh && \
    echo 'echo "Creating configuration file from environment variables..."' >> /app/entrypoint.sh && \
    echo 'cat > /app/settings/conf.yaml << EOF' >> /app/entrypoint.sh && \
    echo 'paths:' >> /app/entrypoint.sh && \
    echo '  logs: "./logs"' >> /app/entrypoint.sh && \
    echo '  assets: "./assets"' >> /app/entrypoint.sh && \
    echo 'database:' >> /app/entrypoint.sh && \
    echo '  dsn: "${DATABASE_DSN}"' >> /app/entrypoint.sh && \
    echo '  debug: ${DATABASE_DEBUG:-true}' >> /app/entrypoint.sh && \
    echo 'scrapping:' >> /app/entrypoint.sh && \
    echo '  total_tries: 3' >> /app/entrypoint.sh && \
    echo '  max_go_rutines: 6' >> /app/entrypoint.sh && \
    echo 'auth:' >> /app/entrypoint.sh && \
    echo '  facebook:' >> /app/entrypoint.sh && \
    echo '    client: "${FACEBOOK_CLIENT_ID:-}"' >> /app/entrypoint.sh && \
    echo '    secret: "${FACEBOOK_CLIENT_SECRET:-}"' >> /app/entrypoint.sh && \
    echo '    callback: "${FACEBOOK_CALLBACK_URL:-}"' >> /app/entrypoint.sh && \
    echo '  google:' >> /app/entrypoint.sh && \
    echo '    client: "${GOOGLE_CLIENT_ID:-}"' >> /app/entrypoint.sh && \
    echo '    secret: "${GOOGLE_CLIENT_SECRET:-}"' >> /app/entrypoint.sh && \
    echo '    callback: "${GOOGLE_CALLBACK_URL:-}"' >> /app/entrypoint.sh && \
    echo 'storage:' >> /app/entrypoint.sh && \
    echo '  path: "./storage"' >> /app/entrypoint.sh && \
    echo 'server:' >> /app/entrypoint.sh && \
    echo '  sk: "${SERVER_SECRET_KEY:-default-secret}"' >> /app/entrypoint.sh && \
    echo '  port: ${PORT:-5000}' >> /app/entrypoint.sh && \
    echo '  host: "0.0.0.0"' >> /app/entrypoint.sh && \
    echo '  debug: false' >> /app/entrypoint.sh && \
    echo 'email:' >> /app/entrypoint.sh && \
    echo '  host: "${EMAIL_HOST:-}"' >> /app/entrypoint.sh && \
    echo '  port: ${EMAIL_PORT:-465}' >> /app/entrypoint.sh && \
    echo '  username: "${EMAIL_USERNAME:-}"' >> /app/entrypoint.sh && \
    echo '  password: "${EMAIL_PASSWORD:-}"' >> /app/entrypoint.sh && \
    echo 'EOF' >> /app/entrypoint.sh && \
    echo 'echo "Configuration file created successfully"' >> /app/entrypoint.sh && \
    echo 'echo "Starting application..."' >> /app/entrypoint.sh && \
    echo 'exec /app/main' >> /app/entrypoint.sh && \
    chmod +x /app/entrypoint.sh

# Expose port
EXPOSE 5000

# Run the application
ENTRYPOINT ["/bin/bash", "/app/entrypoint.sh"]
