$ErrorActionPreference = "Stop"

$staffName = (Read-Host "Staff name").Trim()
$staffEmail = (Read-Host "Staff email").Trim()
$securePassword = Read-Host "Staff password (minimum 8 characters)" -AsSecureString

$passwordPointer = [Runtime.InteropServices.Marshal]::SecureStringToBSTR($securePassword)
try {
    $staffPassword = [Runtime.InteropServices.Marshal]::PtrToStringBSTR($passwordPointer)
}
finally {
    [Runtime.InteropServices.Marshal]::ZeroFreeBSTR($passwordPointer)
}

if ([string]::IsNullOrWhiteSpace($staffName)) {
    throw "Staff name is required."
}

if ([string]::IsNullOrWhiteSpace($staffEmail)) {
    throw "Staff email is required."
}

if ($staffPassword.Length -lt 8) {
    throw "Staff password must contain at least 8 characters."
}

$env:SEED_STAFF_NAME = $staffName
$env:SEED_STAFF_EMAIL = $staffEmail
$env:SEED_STAFF_PASSWORD = $staffPassword

try {
    & (Join-Path $PSScriptRoot "run-local.ps1")
}
finally {
    Remove-Item Env:\SEED_STAFF_NAME -ErrorAction SilentlyContinue
    Remove-Item Env:\SEED_STAFF_EMAIL -ErrorAction SilentlyContinue
    Remove-Item Env:\SEED_STAFF_PASSWORD -ErrorAction SilentlyContinue
    $staffPassword = $null
}
