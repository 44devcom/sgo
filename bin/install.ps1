$ErrorActionPreference = 'Stop'

$AppName = 'sgo'
$Repo = '44devcom/sgo'

Write-Host "Installing $AppName..."

$isWindows = if ($null -ne $IsWindows) { $IsWindows } else { $env:OS -match 'Windows' }
if (-not $isWindows) {
    Write-Error "Unsupported OS: $($env:OS)"
    exit 1
}

$arch = $env:PROCESSOR_ARCHITECTURE
if ($arch -in 'AMD64', 'x86_64') {
    $arch = 'amd64'
} elseif ($arch -in 'ARM64', 'aarch64') {
    $arch = 'arm64'
}

if ($arch -ne 'amd64') {
    Write-Error "Unsupported architecture: $($env:PROCESSOR_ARCHITECTURE). Only windows-amd64 is available. See https://github.com/$Repo#download"
    exit 1
}

$Url = "https://github.com/$Repo/raw/refs/heads/master/dist/windows-amd64/$AppName.exe"
$Dest = Join-Path (Join-Path ([Environment]::GetFolderPath('UserProfile')) 'Downloads') "$AppName.exe"

Write-Host "Downloading from: $Url"

$downloadsDir = Split-Path -Parent $Dest
if (-not (Test-Path $downloadsDir)) {
    New-Item -ItemType Directory -Path $downloadsDir -Force | Out-Null
}

Invoke-WebRequest -Uri $Url -OutFile $Dest -UseBasicParsing

Write-Host ""
Write-Host "Done. Saved to: $Dest"
