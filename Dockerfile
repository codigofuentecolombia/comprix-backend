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

# Copy binary from builder
COPY --from=builder /app/main .

# Create settings directory
RUN mkdir -p settings

# Copy example config as template
COPY --from=builder /app/settings/conf.example.yaml ./settings/

# Create startup script
RUN echo '#!/bin/sh' > /app/start.sh && \
    echo 'cat > /app/settings/conf.yaml << EOF' >> /app/start.sh && \
    echo 'paths:' >> /app/start.sh && \
    echo '  logs: "./logs"' >> /app/start.sh && \
    echo '  assets: "./assets"' >> /app/start.sh && \
    echo 'database:' >> /app/start.sh && \
    echo '  dsn: "${DATABASE_DSN}"' >> /app/start.sh && \
    echo '  debug: ${DATABASE_DEBUG:-true}' >> /app/start.sh && \
    echo 'scrapping:' >> /app/start.sh && \
    echo '  total_tries: 3' >> /app/start.sh && \
    echo '  max_go_rutines: 6' >> /app/start.sh && \
    echo 'auth:' >> /app/start.sh && \
    echo '  facebook:' >> /app/start.sh && \
    echo '    client: "${FACEBOOK_CLIENT_ID}"' >> /app/start.sh && \
    echo '    secret: "${FACEBOOK_CLIENT_SECRET}"' >> /app/start.sh && \
    echo '    callback: "${FACEBOOK_CALLBACK_URL}"' >> /app/start.sh && \
    echo '  google:' >> /app/start.sh && \
    echo '    client: "${GOOGLE_CLIENT_ID}"' >> /app/start.sh && \
    echo '    secret: "${GOOGLE_CLIENT_SECRET}"' >> /app/start.sh && \
    echo '    callback: "${GOOGLE_CALLBACK_URL}"' >> /app/start.sh && \
    echo 'storage:' >> /app/start.sh && \
    echo '  path: "./storage"' >> /app/start.sh && \
    echo 'server:' >> /app/start.sh && \
    echo '  sk: "${SERVER_SECRET_KEY}"' >> /app/start.sh && \
    echo '  port: ${PORT:-5000}' >> /app/start.sh && \
    echo '  host: "0.0.0.0"' >> /app/start.sh && \
    echo '  debug: false' >> /app/start.sh && \
    echo 'email:' >> /app/start.sh && \
    echo '  host: "${EMAIL_HOST}"' >> /app/start.sh && \
    echo '  port: ${EMAIL_PORT:-465}' >> /app/start.sh && \
    echo '  username: "${EMAIL_USERNAME}"' >> /app/start.sh && \
    echo '  password: "${EMAIL_PASSWORD}"' >> /app/start.sh && \
    echo 'EOF' >> /app/start.sh && \
    echo 'exec /app/main' >> /app/start.sh && \
    chmod +x /app/start.sh

# Expose port
EXPOSE 5000

# Run the application
CMD ["/app/start.sh"]
