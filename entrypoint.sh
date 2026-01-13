#!/bin/sh
set -e

# If .env doesn't exist, create it from runtime environment variables
if [ ! -f .env ]; then
  cat > .env <<EOF
DB_HOST=${DB_HOST:-}
DB_PORT=${DB_PORT:-}
DB_NAME=${DB_NAME:-}
DB_USER=${DB_USER:-}
DB_PASS=${DB_PASS:-}
ADMIN_USER=${ADMIN_USER:-}
ADMIN_PASS=${ADMIN_PASS:-}
SECRET=${SECRET:-}
REFRESH_SECRET=${REFRESH_SECRET:-}
DB_SSL_MODE=${DB_SSL_MODE:-disable}
SSL_MODE=${SSL_MODE:-disable}
EOF
fi

exec "$@"
