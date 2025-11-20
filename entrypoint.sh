#!/bin/sh
set -e

echo "Creating configuration file from environment variables..."

cat > /app/settings/conf.yaml << EOF
paths:
  logs: "./logs"
  assets: "./assets"
database:
  dsn: "${DATABASE_DSN}"
  debug: ${DATABASE_DEBUG:-true}
scrapping:
  total_tries: 3
  max_go_rutines: 6
auth:
  facebook:
    client: "${FACEBOOK_CLIENT_ID:-}"
    secret: "${FACEBOOK_CLIENT_SECRET:-}"
    callback: "${FACEBOOK_CALLBACK_URL:-}"
  google:
    client: "${GOOGLE_CLIENT_ID:-}"
    secret: "${GOOGLE_CLIENT_SECRET:-}"
    callback: "${GOOGLE_CALLBACK_URL:-}"
storage:
  path: "./storage"
server:
  sk: "${SERVER_SECRET_KEY:-default-secret}"
  port: ${PORT:-5000}
  host: "0.0.0.0"
  debug: false
email:
  host: "${EMAIL_HOST:-}"
  port: ${EMAIL_PORT:-465}
  username: "${EMAIL_USERNAME:-}"
  password: "${EMAIL_PASSWORD:-}"
EOF

echo "Configuration file created successfully"
echo "Starting application..."

exec /app/main
