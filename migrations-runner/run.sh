#!/bin/sh

set -e

echo "Running pet-service migrations..."
migrate -path /migrations/pet-service -database "$PET_DATABASE_URL" up || true

echo "Running user-service migrations..."
migrate -path /migrations/user-service -database "$USER_DATABASE_URL" up || true

echo "Running order-service migrations..."
migrate -path /migrations/order-service -database "$ORDER_DATABASE_URL" up || true

echo "All migrations completed."
