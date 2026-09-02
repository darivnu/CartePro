#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="$SCRIPT_DIR/../.env"

if [[ ! -f "$ENV_FILE" ]]; then
    echo "Error: $ENV_FILE not found" >&2
    exit 1
fi

set -a
source "$ENV_FILE"
set +a

TABLES="transactions, qr_tokens, clients, partners, admins, employers, sessions, users"

echo "Truncating tables in database '$DB_NAME' on $DB_HOST:$DB_PORT..."

PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" \
    -c "TRUNCATE TABLE $TABLES RESTART IDENTITY CASCADE;"

echo "Done."
