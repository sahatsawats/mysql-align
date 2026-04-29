# Run the mysql-align integration test suite against a temporary MySQL 5.7 container.
# Requires: Docker Desktop running, Go toolchain, port 3307 free (or override $env:MYSQL_PORT).
$ErrorActionPreference = "Stop"

$RepoRoot   = Split-Path -Parent $PSScriptRoot
$ComposeFile = Join-Path $RepoRoot "docker-compose.test.yml"

Write-Host "==> Starting MySQL test container..."
docker compose -f $ComposeFile up -d --wait

try {
    if (-not $env:MYSQL_HOST)     { $env:MYSQL_HOST     = "127.0.0.1" }
    if (-not $env:MYSQL_PORT)     { $env:MYSQL_PORT     = "3307" }
    if (-not $env:MYSQL_USER)     { $env:MYSQL_USER     = "root" }
    if (-not $env:MYSQL_PASSWORD) { $env:MYSQL_PASSWORD = "testpw" }

    Write-Host "==> Running integration tests..."
    Set-Location $RepoRoot
    go test -tags=integration -count=1 -v ./features/...
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
} finally {
    Write-Host "==> Tearing down MySQL test container..."
    docker compose -f $ComposeFile down -v
}
