# PowerShell Script to Update Go on Windows
# This script helps automate the Go update process

param(
    [string]$GoVersion = "1.23.8"
)

Write-Host "=== Go Update Script for Windows ===" -ForegroundColor Green
Write-Host "Target Go Version: $GoVersion" -ForegroundColor Yellow

# Check current Go version
Write-Host "`nChecking current Go version..." -ForegroundColor Cyan
try {
    $currentVersion = go version 2>$null
    if ($currentVersion) {
        Write-Host "Current: $currentVersion" -ForegroundColor White
    } else {
        Write-Host "Go is not currently installed or not in PATH" -ForegroundColor Red
    }
} catch {
    Write-Host "Go is not currently installed or not in PATH" -ForegroundColor Red
}

# Check for vulnerabilities with current version
Write-Host "`nChecking for vulnerabilities..." -ForegroundColor Cyan
try {
    $vulnCheck = govulncheck ./... 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ No vulnerabilities found!" -ForegroundColor Green
        Write-Host "Your Go version is secure. Update may not be necessary." -ForegroundColor Green
        exit 0
    } else {
        Write-Host "🚨 Vulnerabilities detected!" -ForegroundColor Red
        Write-Host "Update is recommended." -ForegroundColor Yellow
    }
} catch {
    Write-Host "Could not check vulnerabilities. Proceeding with update instructions." -ForegroundColor Yellow
}

Write-Host "`n=== Manual Update Instructions ===" -ForegroundColor Green

Write-Host "`n1. Download Go ${GoVersion}:" -ForegroundColor Yellow
Write-Host "   Visit: https://golang.org/dl/" -ForegroundColor White
Write-Host "   Download: go${GoVersion}.windows-amd64.msi" -ForegroundColor White

Write-Host "`n2. Installation Steps:" -ForegroundColor Yellow
Write-Host "   • Run the MSI installer as Administrator" -ForegroundColor White
Write-Host "   • Follow the installation wizard" -ForegroundColor White
Write-Host "   • The installer will automatically update your PATH" -ForegroundColor White

Write-Host "`n3. Verify Installation:" -ForegroundColor Yellow
Write-Host "   • Close and reopen PowerShell" -ForegroundColor White
Write-Host "   • Run: go version" -ForegroundColor White
Write-Host "   • Expected output: go version go${GoVersion} windows/amd64" -ForegroundColor White

Write-Host "`n4. Post-Update Steps:" -ForegroundColor Yellow
Write-Host "   • Run: go clean -cache" -ForegroundColor White
Write-Host "   • Run: go mod tidy" -ForegroundColor White
Write-Host "   • Run: govulncheck ./..." -ForegroundColor White

Write-Host "`n=== Automated Download ===" -ForegroundColor Green
$downloadUrl = "https://golang.org/dl/go${GoVersion}.windows-amd64.msi"
$downloadPath = "$env:TEMP\go${GoVersion}.windows-amd64.msi"

$choice = Read-Host "`nWould you like to download the Go installer automatically? (y/n)"
if ($choice -eq 'y' -or $choice -eq 'Y') {
    try {
        Write-Host "Downloading Go ${GoVersion} installer..." -ForegroundColor Cyan
        Invoke-WebRequest -Uri $downloadUrl -OutFile $downloadPath -UseBasicParsing
        Write-Host "✅ Downloaded to: $downloadPath" -ForegroundColor Green
        
        $runInstaller = Read-Host "Would you like to run the installer now? (y/n)"
        if ($runInstaller -eq 'y' -or $runInstaller -eq 'Y') {
            Write-Host "Starting installer..." -ForegroundColor Cyan
            Start-Process -FilePath $downloadPath -Verb RunAs -Wait
            Write-Host "Installation completed. Please restart PowerShell and run 'go version' to verify." -ForegroundColor Green
        } else {
            Write-Host "Installer saved. Run it manually when ready: $downloadPath" -ForegroundColor Yellow
        }
    } catch {
        Write-Host "❌ Download failed: $($_.Exception.Message)" -ForegroundColor Red
        Write-Host "Please download manually from https://golang.org/dl/" -ForegroundColor Yellow
    }
}

Write-Host "`n=== Security Verification ===" -ForegroundColor Green
Write-Host "After updating, run these commands to verify security:" -ForegroundColor White
Write-Host "  go version" -ForegroundColor Cyan
Write-Host "  govulncheck ./..." -ForegroundColor Cyan

Write-Host "`n=== Documentation ===" -ForegroundColor Green
Write-Host "For detailed security information, see: SECURITY.md" -ForegroundColor White

Write-Host "`nUpdate script completed!" -ForegroundColor Green 