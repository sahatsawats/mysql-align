#!/usr/bin/env bash
# Run the mysql-align integration test suite against a temporary MySQL 5.7 container.
# Requires: Docker Desktop running, Go toolchain, port 3307 free (or override MYSQL_PORT).
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
COMPOSE_FILE="$REPO_ROOT/docker-compose.test.yml"

cleanup() {
    docker compose -f "$COMPOSE_FILE" down -v
}
trap cleanup EXIT

echo "==> Starting MySQL test container..."
docker compose -f "$COMPOSE_FILE" up -d --wait

export MYSQL_HOST="${MYSQL_HOST:-127.0.0.1}"
export MYSQL_PORT="${MYSQL_PORT:-3307}"
export MYSQL_USER="${MYSQL_USER:-root}"
export MYSQL_PASSWORD="${MYSQL_PASSWORD:-testpw}"

echo "==> Running integration tests..."
cd "$REPO_ROOT"
go test -tags=integration -count=1 -v ./features/...
