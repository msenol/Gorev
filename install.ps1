# Gorev Windows Installation Script
# https://github.com/msenol/Gorev

$ErrorActionPreference = "Stop"

# Variables
$REPO = "msenol/Gorev"
$VERSION = if ($env:VERSION) { $env:VERSION } else { "v0.9.0" }
$INSTALL_DIR = if ($env:INSTALL_DIR) { $env:INSTALL_DIR } else { "$env:LOCALAPPDATA\Programs\gorev" }
$DATA_DIR = if ($env:DATA_DIR) { $env:DATA_DIR } else { "$env:APPDATA\gorev" }

Write-Host "Installing Gorev $VERSION..." -ForegroundColor Green
Write-Host "Install directory: $INSTALL_DIR"
Write-Host "Data directory: $DATA_DIR"

# Create directories
New-Item -ItemType Directory -Force -Path $INSTALL_DIR | Out-Null
New-Item -ItemType Directory -Force -Path $DATA_DIR | Out-Null

# Download binary
$BINARY_URL = "https://github.com/$REPO/releases/download/$VERSION/gorev-windows-amd64.exe"
$BINARY_PATH = "$INSTALL_DIR\gorev.exe"

Write-Host "Downloading binary..." -ForegroundColor Yellow
try {
    Invoke-WebRequest -Uri $BINARY_URL -OutFile $BINARY_PATH
} catch {
    Write-Host "Failed to download binary: $_" -ForegroundColor Red
    exit 1
}

# Download source for migrations
$SOURCE_URL = "https://github.com/$REPO/archive/refs/tags/$VERSION.zip"
$TEMP_DIR = New-TemporaryFile | ForEach-Object { Remove-Item $_; New-Item -ItemType Directory -Path $_ }

Write-Host "Downloading data files..." -ForegroundColor Yellow
try {
    $SOURCE_ZIP = "$TEMP_DIR\source.zip"
    Invoke-WebRequest -Uri $SOURCE_URL -OutFile $SOURCE_ZIP
    
    # Extract
    Expand-Archive -Path $SOURCE_ZIP -DestinationPath $TEMP_DIR -Force
    
    # Find source directory
    $VERSION_WITHOUT_V = $VERSION.TrimStart('v')
    $SOURCE_DIR = "$TEMP_DIR\Gorev-$VERSION_WITHOUT_V"
    
    if (!(Test-Path $SOURCE_DIR)) {
        $SOURCE_DIR = Get-ChildItem -Path $TEMP_DIR -Directory | Where-Object { $_.Name -like "Gorev-*" } | Select-Object -First 1
    }
    
    # Copy migrations
    Write-Host "Copying migration files..." -ForegroundColor Yellow
    $MIGRATIONS_SOURCE = "$SOURCE_DIR\gorev-mcpserver\internal\veri\migrations"
    $MIGRATIONS_DEST = "$DATA_DIR\internal\veri\migrations"
    New-Item -ItemType Directory -Force -Path (Split-Path $MIGRATIONS_DEST -Parent) | Out-Null
    Copy-Item -Path $MIGRATIONS_SOURCE -Destination $MIGRATIONS_DEST -Recurse -Force
    
} catch {
    Write-Host "Failed to download or extract data files: $_" -ForegroundColor Red
    exit 1
} finally {
    # Cleanup
    Remove-Item -Path $TEMP_DIR -Recurse -Force -ErrorAction SilentlyContinue
}

# Create batch wrapper
$WRAPPER_PATH = "$INSTALL_DIR\gorev.bat"
$WRAPPER_CONTENT = @"
@echo off
set GOREV_ROOT=$DATA_DIR
"$BINARY_PATH" %*
"@

Set-Content -Path $WRAPPER_PATH -Value $WRAPPER_CONTENT

# Add to PATH if not already there
$PATH = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($PATH -notlike "*$INSTALL_DIR*") {
    Write-Host "Adding to PATH..." -ForegroundColor Yellow
    [Environment]::SetEnvironmentVariable("PATH", "$PATH;$INSTALL_DIR", "User")
    $env:PATH = "$env:PATH;$INSTALL_DIR"
}

# Set GOREV_ROOT permanently
[Environment]::SetEnvironmentVariable("GOREV_ROOT", $DATA_DIR, "User")
$env:GOREV_ROOT = $DATA_DIR

Write-Host "`nGorev installed successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Installed components:"
Write-Host "  Binary: $BINARY_PATH"
Write-Host "  Wrapper: $WRAPPER_PATH"
Write-Host "  Data files: $DATA_DIR"
Write-Host ""
Write-Host "Run 'gorev version' to verify installation"
Write-Host "Run 'gorev help' to see available commands"
Write-Host ""
Write-Host "NOTE: You may need to restart your terminal for PATH changes to take effect."