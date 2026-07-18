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
  $startInfo = New-Object System.Diagnostics.ProcessStartInfo
  $startInfo.FileName = $Path
  $startInfo.Arguments = '--version'
  $startInfo.UseShellExecute = $false
  $startInfo.CreateNoWindow = $true
  $startInfo.RedirectStandardOutput = $true
  $startInfo.RedirectStandardError = $true
  $process = New-Object System.Diagnostics.Process
  $process.StartInfo = $startInfo
  try {
    if (-not $process.Start()) { throw 'Could not start installed RunX for version verification' }
    if (-not $process.WaitForExit(10000)) {
      $process.Kill()
      $process.WaitForExit()
      throw 'Installed RunX version check timed out after 10 seconds'
    }
    $stdout = $process.StandardOutput.ReadToEnd().Trim()
    $stderr = $process.StandardError.ReadToEnd().Trim()
    if ($process.ExitCode -ne 0) { throw "Installed RunX exited with code $($process.ExitCode) during verification: $stderr" }
    if ($stdout -ne $ExpectedVersion) { throw "Installed RunX reported $stdout; expected $ExpectedVersion" }
  } finally {
    $process.Dispose()
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
  $firstAsset = $assetCandidates[0]
  $sourceUrl = if ($DownloadBaseUrl) { "$DownloadBaseUrl/$encodedTag/$firstAsset" } else { "https://github.com/$Repo/releases/download/$encodedTag/$firstAsset" }
  Write-Host 'Initiating GUIHO CLI Upgrade / Installation Sequence...'
  Write-Host "Target Version: v$targetVersion"
  Write-Host "Architecture:   $detectedArch"
  Write-Host "Variant:        $variant"
  Write-Host "Source URL:     $sourceUrl"
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
  $skillAsset = Join-Path $temporaryDirectory 'guiho-s-runx'
  $promptAsset = Join-Path $temporaryDirectory 'guiho-i-runx'
  $assetBase = if ($DownloadBaseUrl) { "$DownloadBaseUrl/$encodedTag" } else { "https://github.com/$Repo/releases/download/$encodedTag" }
  Write-Host "Downloading skill asset: $assetBase/guiho-s-runx"
  Invoke-WebRequest -Uri "$assetBase/guiho-s-runx" -OutFile $skillAsset -UseBasicParsing
  Write-Host "Downloading instruction asset: $assetBase/guiho-i-runx"
  Invoke-WebRequest -Uri "$assetBase/guiho-i-runx" -OutFile $promptAsset -UseBasicParsing
  foreach ($skillRoot in @((Join-Path $HOME '.agents\skills\guiho-s-runx'), (Join-Path $HOME '.claude\skills\guiho-s-runx'))) {
    New-Item -ItemType Directory -Force -Path $skillRoot | Out-Null
    Copy-Item -LiteralPath $skillAsset -Destination (Join-Path $skillRoot 'SKILL.md') -Force
    Write-Host "Installed skill: $(Join-Path $skillRoot 'SKILL.md')"
  }
  $instructionTargets = @()
  if (Test-Path -LiteralPath (Join-Path (Get-Location) 'AGENTS.md')) { $instructionTargets += (Join-Path (Get-Location) 'AGENTS.md') }
  if (Test-Path -LiteralPath (Join-Path (Get-Location) 'CLAUDE.md')) { $instructionTargets += (Join-Path (Get-Location) 'CLAUDE.md') }
  if ($instructionTargets.Count -eq 0) { $instructionTargets += (Join-Path (Get-Location) 'AGENTS.md') }
  $startMarker = '<!-- BEGIN RUNX — DO NOT EDIT THIS SECTION -->'
  $endMarker = '<!-- END RUNX -->'
  $prompt = Get-Content -Raw -LiteralPath $promptAsset
  foreach ($instructionPath in $instructionTargets) {
    Write-Host "Reconciling instruction file: $instructionPath"
    $existing = if (Test-Path -LiteralPath $instructionPath) { Get-Content -Raw -LiteralPath $instructionPath } else { '' }
    $pattern = [Regex]::Escape($startMarker) + '[\s\S]*?' + [Regex]::Escape($endMarker) + '\s*'
    $clean = ([Regex]::Replace($existing, $pattern, '')).TrimEnd()
    $prefix = if ($clean) { "$clean`r`n`r`n" } else { '' }
    Set-Content -LiteralPath $instructionPath -Value "$prefix$startMarker`r`n$($prompt.Trim())`r`n$endMarker`r`n" -Encoding utf8
  }
  if ($env:RUNX_SKIP_PATH_UPDATE -ne '1') {
    Add-InstallDirToPath -Directory $InstallDir
    Test-Shadowing -ExpectedPath $destination
  }
  Write-Host "Final verification: $destination --version"
  Write-Host "Installed and verified RunX $targetVersion at $destination"
} finally {
  Remove-Item -LiteralPath $temporaryDirectory -Recurse -Force -ErrorAction SilentlyContinue
}
