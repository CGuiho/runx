param(
  [string]$Version,
  [string]$Arch,
  [string]$Variant,
  [string]$InstallDir,
  [switch]$Help
)

$ErrorActionPreference = 'Stop'
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

if ([string]::IsNullOrWhiteSpace($Version)) { $Version = if ($env:RUNX_VERSION) { $env:RUNX_VERSION } else { 'latest' } }
$Repo = if ($env:RUNX_REPO) { $env:RUNX_REPO } else { 'CGuiho/runx' }
$DownloadBaseUrl = if ($env:RUNX_DOWNLOAD_BASE_URL) { $env:RUNX_DOWNLOAD_BASE_URL.TrimEnd('/') } else { $null }
if ([string]::IsNullOrWhiteSpace($InstallDir)) { $InstallDir = if ($env:RUNX_INSTALL_DIR) { $env:RUNX_INSTALL_DIR } else { Join-Path $HOME '.local\bin' } }

if ($Help -or $Version -eq '--help' -or $Version -eq '-h') {
  @"
Install GUIHO RunX as a verified native CLI binary from GitHub Releases.

Usage: install.ps1 [-Version VERSION] [-Arch ARCH] [-Variant VARIANT] [-InstallDir DIR]

Parameters:
  -Version      Exact stable or prerelease version (default: latest stable).
  -Arch         Force architecture: x64 | arm64 (default: auto-detect).
  -Variant      Force x64 variant: baseline | default | modern (default: baseline).
  -InstallDir   Install directory (default: `$HOME\.local\bin).
  -Help         Show this help.
"@
  return
}

function Resolve-TargetVersion {
  param([string]$RequestedVersion)
  $tag = $RequestedVersion
  if ($RequestedVersion -eq 'latest') {
    Write-Host 'Resolving latest stable RunX release...'
    $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -Headers @{ Accept = 'application/vnd.github+json' }
    $tag = [string]$release.tag_name
  }
  $normalized = $tag -replace '^@guiho/runx@', '' -replace '^v', ''
  if ($normalized -notmatch '^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)?(?:\+[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)?$') {
    throw "Invalid RunX version: $RequestedVersion"
  }
  $prerelease = (($normalized -split '\+', 2)[0] -split '-', 2)
  if ($prerelease.Count -eq 2) {
    foreach ($identifier in $prerelease[1] -split '\.') {
      if ($identifier -match '^0\d+$') { throw "Invalid RunX version: $RequestedVersion" }
    }
  }
  return $normalized
}

function Get-PathEntries { param([string]$PathValue); if ([string]::IsNullOrWhiteSpace($PathValue)) { return @() }; return @($PathValue -split ';' | Where-Object { -not [string]::IsNullOrWhiteSpace($_) }) }
function Test-PathContains {
  param([string]$PathValue, [string]$Directory)
  $normalizedDirectory = $Directory.TrimEnd('\')
  foreach ($entry in Get-PathEntries -PathValue $PathValue) {
    if ($entry.TrimEnd('\').Equals($normalizedDirectory, [StringComparison]::OrdinalIgnoreCase)) { return $true }
  }
  return $false
}
function Add-InstallDirToPath {
  param([string]$Directory)
  $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
  if (-not (Test-PathContains -PathValue $userPath -Directory $Directory)) {
    $newUserPath = (@($Directory) + @(Get-PathEntries -PathValue $userPath)) -join ';'
    [Environment]::SetEnvironmentVariable('Path', $newUserPath.TrimEnd(';'), 'User')
    Write-Host "Added $Directory to user PATH. Restart your terminal to use runx globally."
  }
  if (-not (Test-PathContains -PathValue $env:Path -Directory $Directory)) { $env:Path = "$Directory;$env:Path" }
}

function Test-NativeBinary {
  param([string]$Path)
  $stream = [System.IO.File]::OpenRead($Path)
  try { return $stream.Length -ge 2 -and $stream.ReadByte() -eq 0x4D -and $stream.ReadByte() -eq 0x5A }
  finally { $stream.Dispose() }
}

function Test-InstalledVersion {
  param([string]$Path, [string]$ExpectedVersion)
  $stdoutPath = [System.IO.Path]::GetTempFileName()
  $stderrPath = [System.IO.Path]::GetTempFileName()
  try {
    $process = Start-Process -FilePath $Path -ArgumentList '--version' -PassThru -NoNewWindow -RedirectStandardOutput $stdoutPath -RedirectStandardError $stderrPath
    if (-not $process.WaitForExit(10000)) {
      $process.Kill()
      $process.WaitForExit()
      throw 'Installed RunX version check timed out after 10 seconds'
    }
    $process.WaitForExit()
    $stdout = ([string](Get-Content -LiteralPath $stdoutPath -Raw -ErrorAction SilentlyContinue)).Trim()
    $stderr = ([string](Get-Content -LiteralPath $stderrPath -Raw -ErrorAction SilentlyContinue)).Trim()
    if ($process.ExitCode -ne 0) { throw "Installed RunX exited with code $($process.ExitCode) during verification: $stderr" }
    if ($stdout -ne $ExpectedVersion) { throw "Installed RunX reported $stdout; expected $ExpectedVersion" }
  } finally {
    Remove-Item -LiteralPath $stdoutPath, $stderrPath -Force -ErrorAction SilentlyContinue
  }
}

function Start-BackupCleanup {
  param([string]$BackupPath)
  $script = 'for ($attempt = 0; $attempt -lt 300; $attempt += 1) { try { Remove-Item -LiteralPath $env:RUNX_BACKUP_PATH -Force -ErrorAction Stop; exit 0 } catch { Start-Sleep -Milliseconds 100 } }; exit 1'
  $encoded = [Convert]::ToBase64String([Text.Encoding]::Unicode.GetBytes($script))
  $previousBackupPath = $env:RUNX_BACKUP_PATH
  try {
    $env:RUNX_BACKUP_PATH = $BackupPath
    Start-Process powershell.exe -ArgumentList @('-NoLogo', '-NoProfile', '-NonInteractive', '-EncodedCommand', $encoded) -WindowStyle Hidden
  } finally {
    $env:RUNX_BACKUP_PATH = $previousBackupPath
  }
}

function Install-Transactional {
  param([string]$DownloadedPath, [string]$Destination, [string]$ExpectedVersion)
  $backupPath = "$Destination.old-$PID-$([Guid]::NewGuid().ToString('N'))"
  $originalMoved = $false
  try {
    if (Test-Path -LiteralPath $Destination) { Move-Item -LiteralPath $Destination -Destination $backupPath; $originalMoved = $true }
    Move-Item -LiteralPath $DownloadedPath -Destination $Destination
    Write-Host 'Verifying...'
    Test-InstalledVersion -Path $Destination -ExpectedVersion $ExpectedVersion
    if ($originalMoved) {
      try { Remove-Item -LiteralPath $backupPath -Force }
      catch { Start-BackupCleanup -BackupPath $backupPath; Write-Host 'Old executable cleanup will finish after the running process exits.' }
    }
  } catch {
    $failure = $_.Exception.Message
    try {
      if (Test-Path -LiteralPath $Destination) { Remove-Item -LiteralPath $Destination -Force -ErrorAction Stop }
      if ($originalMoved) {
        if (-not (Test-Path -LiteralPath $backupPath)) { throw "Backup is missing: $backupPath" }
        Move-Item -LiteralPath $backupPath -Destination $Destination -ErrorAction Stop
      }
    } catch {
      throw "RunX installation failed: $failure. Automatic rollback also failed: $($_.Exception.Message). Backup remains at $backupPath"
    }
    if ($originalMoved) { throw "RunX installation failed and the previous executable was restored: $failure" }
    throw "RunX installation failed before a previous executable could be restored: $failure"
  }
}

function Test-Shadowing {
  param([string]$ExpectedPath)
  $command = Get-Command runx -ErrorAction SilentlyContinue
  if ($command -and -not $command.Source.Equals($ExpectedPath, [StringComparison]::OrdinalIgnoreCase)) {
    Write-Warning "Another runx appears earlier in PATH: $($command.Source)"
    Write-Warning "The newly installed binary is at: $ExpectedPath"
  }
}

$detectedArch = if ($Arch) { $Arch } else { switch ($env:PROCESSOR_ARCHITECTURE) { 'AMD64' { 'x64' } 'ARM64' { 'arm64' } default { throw "Unsupported architecture: $env:PROCESSOR_ARCHITECTURE" } } }
if ($detectedArch -notin @('x64', 'arm64')) { throw "Invalid architecture: $detectedArch" }
if (-not [Environment]::Is64BitOperatingSystem) { throw 'Unsupported platform: Windows 32-bit is not supported.' }
$variant = if ($Variant) { $Variant } else { 'baseline' }
$assetCandidates = if ($detectedArch -eq 'arm64') {
  if ($Variant) { throw '-Variant is only valid for x64 installs.' }
  @('runx-windows-arm64.exe')
} else {
  switch ($variant) {
    'baseline' { @('runx-windows-x64-baseline.exe', 'runx-windows-x64.exe', 'runx-windows-x64-modern.exe') }
    'default' { @('runx-windows-x64.exe', 'runx-windows-x64-baseline.exe', 'runx-windows-x64-modern.exe') }
    'modern' { @('runx-windows-x64-modern.exe', 'runx-windows-x64.exe', 'runx-windows-x64-baseline.exe') }
    default { throw "Invalid variant: $variant" }
  }
}

$targetVersion = Resolve-TargetVersion -RequestedVersion $Version
$encodedTag = [Uri]::EscapeDataString("@guiho/runx@$targetVersion")
$temporaryDirectory = Join-Path ([System.IO.Path]::GetTempPath()) ("runx-install-$PID-$([Guid]::NewGuid().ToString('N'))")
New-Item -ItemType Directory -Force -Path $temporaryDirectory | Out-Null

try {
  New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
  $InstallDir = (Resolve-Path -LiteralPath $InstallDir).Path
  $destination = Join-Path $InstallDir 'runx.exe'
  Write-Host "Installing RunX $targetVersion  os=windows  arch=$detectedArch  path=$destination"
  $downloadedPath = $null
  foreach ($asset in $assetCandidates) {
    $url = if ($DownloadBaseUrl) { "$DownloadBaseUrl/$encodedTag/$asset" } else { "https://github.com/$Repo/releases/download/$encodedTag/$asset" }
    $candidatePath = Join-Path $temporaryDirectory $asset
    Write-Host "Downloading $url"
    try {
      Invoke-WebRequest -Uri $url -OutFile $candidatePath -UseBasicParsing
    } catch {
      $statusCode = if ($_.Exception.Response -and $_.Exception.Response.StatusCode) { [int]$_.Exception.Response.StatusCode } else { 0 }
      if ($statusCode -eq 404) { Write-Host "  $asset is not available; trying the next compatible candidate."; continue }
      throw "Failed to download $url`: $($_.Exception.Message)"
    }
    if (-not (Test-NativeBinary $candidatePath)) { throw "Downloaded asset $asset is not a native Windows executable." }
    Unblock-File -LiteralPath $candidatePath -ErrorAction SilentlyContinue
    $downloadedPath = $candidatePath
    break
  }
  if (-not $downloadedPath) { throw "No compatible RunX $targetVersion binary found at https://github.com/$Repo/releases" }
  Write-Host 'Replacing...'
  Install-Transactional -DownloadedPath $downloadedPath -Destination $destination -ExpectedVersion $targetVersion
  if ($env:RUNX_SKIP_PATH_UPDATE -ne '1') {
    Add-InstallDirToPath -Directory $InstallDir
    Test-Shadowing -ExpectedPath $destination
  }
  Write-Host "Installed and verified RunX $targetVersion at $destination"
} finally {
  Remove-Item -LiteralPath $temporaryDirectory -Recurse -Force -ErrorAction SilentlyContinue
}
