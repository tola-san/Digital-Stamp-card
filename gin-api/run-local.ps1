$ErrorActionPreference = "Stop"

$env:DATABASE_URL = "mysql://app:app@127.0.0.1:3306/digital_stamp?ssl-mode=DISABLED"
$env:FRONTEND_URL = "http://localhost:3000"
$env:COOKIE_SECURE = "false"
$env:PORT = "8080"
$env:GOCACHE = Join-Path $PSScriptRoot ".gocache"

Set-Location -LiteralPath $PSScriptRoot
go run ./cmd/api
