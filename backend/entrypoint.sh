#!/bin/sh

set -e

echo "Running database migrations..."
goose up -dir "./database/migrations"
echo "Migrations completed."

if [ "${GO_PROD:-0}" -eq 1 ]; then
    echo "Starting in PRODUCTION mode..."
    exec /tmp/main
else
    echo "Starting in DEVELOPMENT mode with Air..."
    exec air -c .air.toml
fi