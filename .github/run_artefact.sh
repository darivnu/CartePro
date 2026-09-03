#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")"

[ -f .env ] || cp env.example .env

echo "Loading pre-built images..."
docker load -i images.tar

echo "Starting CartePro (frontend, backend, database)..."
docker compose up -d

cat <<'EOF'

CartePro is up:
  Frontend: http://localhost:3000
  Backend:  http://localhost:4242

Logs: docker compose logs -f
Stop: docker compose down
EOF
