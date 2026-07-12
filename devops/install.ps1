param(
  [string]$Version,
  [string]$Arch,
  [string]$Variant,
  [string]$InstallDir,
  [switch]$Help
)

$ErrorActionPreference = 'Stop'
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

# === Defaults from env vars or sensible defaults ===
if ([string]::IsNullOrWhiteSpace($Version)) {
  $Version = if ($env:RUNX_VERSION) { $env:RUNX_VERSION } else { 'latest' }
}
$Repo = if ($env:RUNX_REPO) { $env:RUNX_REPO } else { 'CGuiho/runx' }
if ([string]::IsNullOrWhiteSpace($InstallDir)) {
  $InstallDir = if ($env:RUNX_INSTALL_DIR) { $env:RUNX_INSTALL_DIR } else { Join-Path $HOME '.local\bin' }
}

# === Show help ===
if ($Help -or $Version -eq '--help' -or $Version -eq '-h') {
  @"
Install GUIHO RunX as a native CLI binary from GitHub Releases.

Usage: install.ps1 [-Version VERSION] [-Arch ARCH] [-Variant VARIANT] [-InstallDir DIR]

Parameters:
  -Version      Version to install (default: latest).
                Examples: latest, 0.1.1, @guiho/runx@0.1.1
  -Arch         Force architecture: x64 | arm64 (default: auto-detect)
  -Variant      Force x64 variant: baseline | default | modern (default: baseline)
  -InstallDir   Install directory (default: `$HOME\.local\bin)
  -Help         Show this help

Environment variables:
  RUNX_VERSION, RUNX_REPO, RUNX_INSTALL_DIR
"@
  return
}

# === Detect architecture (compatible with PowerShell 5.1+) ===
$detectedArch = if ($Arch) {
  $Arch
} else {
  switch ($env:PROCESSOR_ARCHITECTURE) {
    'AMD64' { 'x64' }
    'ARM64' { 'arm64' }
    default { throw "Unsupported architecture: $env:PROCESSOR_ARCHITECTURE. Must be AMD64 or ARM64." }
  }
}

if ($detectedArch -notin @('x64', 'arm64')) {
  throw "Invalid architecture: $detectedArch. Must be x64 or arm64."
}

if (-not [Environment]::Is64BitOperatingSystem) {
  throw 'Unsupported platform: Windows 32-bit is not supported.'
}

# === Build asset candidates (baseline-first for x64) ===
$variant = if ($Variant) { $Variant } else { 'baseline' }

$assetCandidates = if ($detectedArch -eq 'x64') {
  switch ($variant) {
    'baseline' { @(
      "runx-windows-x64-baseline.exe",
      "runx-windows-x64.exe",
      "runx-windows-x64-modern.exe"
    )}
    'default' { @(
      "runx-windows-x64.exe",
      "runx-windows-x64-baseline.exe",
      "runx-windows-x64-modern.exe"
    )}
    'modern' { @(
      "runx-windows-x64-modern.exe",
      "runx-windows-x64.exe",
      "runx-windows-x64-baseline.exe"
    )}
    default { throw "Invalid variant: $variant. Must be baseline, default, or modern." }
  }
} else {
  if ($Variant) {
    throw '-Variant is only valid for x64 installs.'
  }
  @("runx-windows-arm64.exe")
}

# === Build download URL ===
function Get-DownloadUrl {
  param([string]$Asset)

  if ($Version -eq 'latest') {
    return "https://github.com/$Repo/releases/latest/download/$Asset"
  }

  $tag = if ($Version.StartsWith('@guiho/runx@')) { $Version }
         elseif ($Version.StartsWith('@')) { $Version }
         else { "@guiho/runx@$Version" }

  $encodedTag = [Uri]::EscapeDataString($tag)
  return "https://github.com/$Repo/releases/download/$encodedTag/$Asset"
}

function Get-PathEntries {
  param([string]$PathValue)

  if ([string]::IsNullOrWhiteSpace($PathValue)) {
    return @()
  }

  return @($PathValue -split ';' | Where-Object { -not [string]::IsNullOrWhiteSpace($_) })
}

function Test-PathContains {
  param(
    [string]$PathValue,
    [string]$Directory
  )

  $normalizedDirectory = $Directory.TrimEnd('\')
  foreach ($entry in Get-PathEntries -PathValue $PathValue) {
    if ($entry.TrimEnd('\').Equals($normalizedDirectory, [StringComparison]::OrdinalIgnoreCase)) {
      return $true
    }
  }

  return $false
}

function Add-InstallDirToPath {
  param([string]$Directory)

  $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
  if (-not (Test-PathContains -PathValue $userPath -Directory $Directory)) {
    $entries = @(Get-PathEntries -PathValue $userPath)
    $newUserPath = (@($Directory) + $entries) -join ';'
    [Environment]::SetEnvironmentVariable('Path', $newUserPath.TrimEnd(';'), 'User')
    Write-Host "Added $Directory to user PATH. Restart your terminal to use runx globally."
  } else {
    Write-Host "$Directory is already configured in user PATH."
  }

  if (-not (Test-PathContains -PathValue $env:Path -Directory $Directory)) {
    $env:Path = "$Directory;$env:Path"
  }
}

$temporaryFile = Join-Path ([System.IO.Path]::GetTempPath()) ([System.IO.Path]::GetRandomFileName())

function Test-NativeBinary {
  param([string]$Path)

  $bytes = [System.IO.File]::ReadAllBytes($Path)
  if ($bytes.Length -lt 2) {
    return $false
  }

  return $bytes[0] -eq 0x4D -and $bytes[1] -eq 0x5A
}

function Test-Shadowing {
  param([string]$ExpectedPath)

  $command = Get-Command runx -ErrorAction SilentlyContinue
  if (-not $command) {
    return
  }

  if (-not $command.Source.Equals($ExpectedPath, [StringComparison]::OrdinalIgnoreCase)) {
    Write-Warning "Another runx appears earlier in PATH: $($command.Source)"
    Write-Warning "The newly installed binary is at: $ExpectedPath"
  }
}

# === Main ===
$variantLabel = if ($Variant) { " variant=$Variant" } else { "" }
Write-Host "runx: $Version  os=windows  arch=$detectedArch$variantLabel"

New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
$InstallDir = (Resolve-Path -LiteralPath $InstallDir).Path
$destination = Join-Path $InstallDir 'runx.exe'

foreach ($asset in $assetCandidates) {
  $url = Get-DownloadUrl -Asset $asset
  Write-Host "  Trying $url"
  try {
    Invoke-WebRequest -Uri $url -OutFile $temporaryFile -UseBasicParsing -ErrorAction Stop
    if (-not (Test-NativeBinary -Path $temporaryFile)) {
      Write-Host "  $asset was not a native Windows binary, trying next..."
      continue
    }

    Move-Item -Force -Path $temporaryFile -Destination $destination
    Write-Host "Installed runx to $destination"

    Add-InstallDirToPath -Directory $InstallDir
    Test-Shadowing -ExpectedPath $destination

    if (Test-Path -LiteralPath $temporaryFile) {
      Remove-Item -LiteralPath $temporaryFile -Force
    }

    Write-Host 'Run: runx --help'
    return
  } catch {
    Write-Host "  not available, trying next..."
  }
}

if (Test-Path -LiteralPath $temporaryFile) {
  Remove-Item -LiteralPath $temporaryFile -Force
}

throw "No compatible runx binary found. Check available assets at: https://github.com/$Repo/releases"
