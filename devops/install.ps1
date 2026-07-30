param(
  [string]$Version = $(if ($env:RUNX_VERSION) { $env:RUNX_VERSION } else { 'latest' }),
  [string]$InstallDir = $(if ($env:RUNX_INSTALL_DIR) { $env:RUNX_INSTALL_DIR } else { Join-Path $HOME '.local\bin' }),
  [switch]$Help
)

$ErrorActionPreference = 'Stop'

function Get-ComparablePath {
  param([Parameter(Mandatory = $true)][string]$Value)

  $Expanded = [Environment]::ExpandEnvironmentVariables($Value.Trim().Trim('"'))
  try {
    return [IO.Path]::GetFullPath($Expanded).TrimEnd([char[]]@('\', '/'))
  } catch {
    return $Expanded.TrimEnd([char[]]@('\', '/'))
  }
}

function Test-PathValueContains {
  param(
    [AllowNull()][string]$PathValue,
    [Parameter(Mandatory = $true)][string]$Directory
  )

  if ([string]::IsNullOrWhiteSpace($PathValue)) { return $false }
  $Expected = Get-ComparablePath -Value $Directory
  foreach ($Entry in ($PathValue -split ';')) {
    if ([string]::IsNullOrWhiteSpace($Entry)) { continue }
    $Actual = Get-ComparablePath -Value $Entry
    if ([string]::Equals($Actual, $Expected, [StringComparison]::OrdinalIgnoreCase)) {
      return $true
    }
  }
  return $false
}

function Add-WindowsPath {
  param([Parameter(Mandatory = $true)][string]$InstallDirectory)

  $ResolvedInstallDirectory = Get-ComparablePath -Value $InstallDirectory
  $UserPath = [Environment]::GetEnvironmentVariable('Path', 'User')
  $UserPathChanged = $false
  if (-not (Test-PathValueContains -PathValue $UserPath -Directory $ResolvedInstallDirectory)) {
    $UpdatedUserPath = if ([string]::IsNullOrWhiteSpace($UserPath)) {
      $ResolvedInstallDirectory
    } else {
      "$ResolvedInstallDirectory;$UserPath"
    }
    [Environment]::SetEnvironmentVariable('Path', $UpdatedUserPath, 'User')
    $UserPathChanged = $true
  }

  $PersistedUserPath = [Environment]::GetEnvironmentVariable('Path', 'User')
  if (-not (Test-PathValueContains -PathValue $PersistedUserPath -Directory $ResolvedInstallDirectory)) {
    throw "Failed to add $ResolvedInstallDirectory to the Windows user Path."
  }

  $ProcessPathChanged = $false
  if (-not (Test-PathValueContains -PathValue $env:Path -Directory $ResolvedInstallDirectory)) {
    $env:Path = if ([string]::IsNullOrWhiteSpace($env:Path)) {
      $ResolvedInstallDirectory
    } else {
      "$ResolvedInstallDirectory;$env:Path"
    }
    $ProcessPathChanged = $true
  }

  return [pscustomobject]@{
    Directory = $ResolvedInstallDirectory
    UserPathChanged = $UserPathChanged
    ProcessPathChanged = $ProcessPathChanged
  }
}

function ConvertTo-GitBashPath {
  param(
    [Parameter(Mandatory = $true)][string]$InstallDirectory,
    [Parameter(Mandatory = $true)][string]$HomeDirectory
  )

  $ResolvedInstallDirectory = Get-ComparablePath -Value $InstallDirectory
  $ResolvedHomeDirectory = Get-ComparablePath -Value $HomeDirectory
  $HomePrefix = $ResolvedHomeDirectory + [IO.Path]::DirectorySeparatorChar
  if ($ResolvedInstallDirectory.StartsWith($HomePrefix, [StringComparison]::OrdinalIgnoreCase)) {
    $Relative = $ResolvedInstallDirectory.Substring($HomePrefix.Length).Replace('\', '/')
    return '$HOME/' + $Relative
  }

  if ($ResolvedInstallDirectory -match '^([A-Za-z]):\\(.*)$') {
    return '/' + $Matches[1].ToLowerInvariant() + '/' + $Matches[2].Replace('\', '/')
  }
  return $ResolvedInstallDirectory.Replace('\', '/')
}

function Add-GitBashPath {
  param(
    [Parameter(Mandatory = $true)][string]$InstallDirectory,
    [string]$HomeDirectory = $HOME
  )

  if ([string]::IsNullOrWhiteSpace($HomeDirectory)) {
    throw 'Cannot configure Git Bash because the home directory is empty.'
  }

  $GitBashPath = ConvertTo-GitBashPath -InstallDirectory $InstallDirectory -HomeDirectory $HomeDirectory
  if ($GitBashPath.StartsWith('$HOME/')) {
    $ExportLine = 'export PATH="' + $GitBashPath + ':$PATH"'
  } else {
    $EscapedPath = $GitBashPath.Replace('\', '\\').Replace('"', '\"').Replace('$', '\$').Replace('`', '\`')
    $ExportLine = 'export PATH="' + $EscapedPath + ':$PATH"'
  }

  $ProfilePath = Join-Path $HomeDirectory '.bashrc'
  $Existing = if (Test-Path -LiteralPath $ProfilePath) {
    [IO.File]::ReadAllText($ProfilePath)
  } else {
    ''
  }
  $LinePattern = '(?m)^[ \t]*' + [Regex]::Escape($ExportLine) + '[ \t]*\r?$'
  $Changed = -not [Regex]::IsMatch($Existing, $LinePattern)
  if ($Changed) {
    $Newline = if ($Existing.Contains("`r`n")) { "`r`n" } else { "`n" }
    $Prefix = if ($Existing.Length -eq 0 -or $Existing.EndsWith("`n") -or $Existing.EndsWith("`r")) { '' } else { $Newline }
    $Addition = $Prefix + '# Added by the RunX installer.' + $Newline + $ExportLine + $Newline
    $Utf8NoBom = New-Object System.Text.UTF8Encoding -ArgumentList $false
    [IO.File]::AppendAllText($ProfilePath, $Addition, $Utf8NoBom)
  }

  return [pscustomobject]@{
    ProfilePath = $ProfilePath
    ExportLine = $ExportLine
    Changed = $Changed
  }
}

if ($Help) { Write-Output 'Usage: install.ps1 [-Version VERSION] [-InstallDir DIRECTORY]'; return }
if ($env:RUNX_INSTALLER_SOURCE_ONLY -eq '1') { return }
$Repo = if ($env:RUNX_REPO) { $env:RUNX_REPO } else { 'CGuiho/runx' }
$Asset = switch ($env:PROCESSOR_ARCHITECTURE.ToUpperInvariant()) { 'AMD64' { 'runx-windows-amd64.exe' } 'ARM64' { 'runx-windows-arm64.exe' } default { throw "Unsupported Windows architecture: $env:PROCESSOR_ARCHITECTURE" } }
if ($Version -eq 'latest') { $BaseUrl = "https://github.com/$Repo/releases/latest/download"; $TargetLabel = 'latest' }
else {
  $Version = $Version -replace '^@guiho/runx/v','' -replace '^@guiho/runx@','' -replace '^v',''
  if ($Version -notmatch '^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?$') { throw "Invalid RunX version: $Version" }
  $EncodedTag = [Uri]::EscapeDataString("@guiho/runx/v$Version"); $BaseUrl = "https://github.com/$Repo/releases/download/$EncodedTag"; $TargetLabel = $Version
}
if ($env:RUNX_DOWNLOAD_BASE_URL) { $BaseUrl = $env:RUNX_DOWNLOAD_BASE_URL.TrimEnd('/') }
$TempDir = Join-Path ([IO.Path]::GetTempPath()) ("runx-install-" + [Guid]::NewGuid().ToString('N')); New-Item -ItemType Directory -Path $TempDir | Out-Null
try {
  Write-Host 'Initiating GUIHO CLI Upgrade / Installation Sequence...'; Write-Host "Target Version: $TargetLabel"; Write-Host "Architecture:   $env:PROCESSOR_ARCHITECTURE"; Write-Host "Target Asset:   $Asset"; Write-Host "Source URL:     $BaseUrl/$Asset"
  $Binary = Join-Path $TempDir $Asset; $Checksums = Join-Path $TempDir 'checksums.txt'; $Skill = Join-Path $TempDir 'guiho-s-runx.zip'
  Invoke-WebRequest -Uri "$BaseUrl/$Asset" -OutFile $Binary -UseBasicParsing; Invoke-WebRequest -Uri "$BaseUrl/checksums.txt" -OutFile $Checksums -UseBasicParsing; Invoke-WebRequest -Uri "$BaseUrl/guiho-s-runx.zip" -OutFile $Skill -UseBasicParsing
  foreach ($Item in @(@($Binary,$Asset),@($Skill,'guiho-s-runx.zip'))) { $Line = Get-Content -LiteralPath $Checksums | Where-Object { $_ -match "\s+$([Regex]::Escape($Item[1]))$" } | Select-Object -First 1; if (-not $Line) { throw "Checksum entry missing for $($Item[1])" }; $Expected = ($Line -split '\s+')[0]; $Actual = (Get-FileHash -Algorithm SHA256 -LiteralPath $Item[0]).Hash; if ($Expected -ne $Actual) { throw "Checksum verification failed for $($Item[1])" } }
  Write-Host '[OK] SHA-256 verification complete.'; New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null; $Destination = Join-Path $InstallDir 'runx.exe'; $Backup = "$Destination.old-$PID"; $HadBackup = $false
  try { if (Test-Path -LiteralPath $Destination) { Move-Item -LiteralPath $Destination -Destination $Backup; $HadBackup = $true }; Move-Item -LiteralPath $Binary -Destination $Destination; $PreviousUpdate=$env:RUNX_DISABLE_UPDATE_WORKER; $PreviousMaintenance=$env:RUNX_DISABLE_AGENT_MAINTENANCE_WORKER; $env:RUNX_DISABLE_UPDATE_WORKER='1';$env:RUNX_DISABLE_AGENT_MAINTENANCE_WORKER='1'; try { $ActualVersion=(& $Destination --version).Trim() } finally {$env:RUNX_DISABLE_UPDATE_WORKER=$PreviousUpdate;$env:RUNX_DISABLE_AGENT_MAINTENANCE_WORKER=$PreviousMaintenance}; if ($Version -ne 'latest' -and $ActualVersion -ne $Version) { throw "Installed version $ActualVersion does not match $Version" }; if ($HadBackup) { Remove-Item -LiteralPath $Backup -Force } } catch { if (Test-Path -LiteralPath $Destination) { Remove-Item -LiteralPath $Destination -Force }; if ($HadBackup -and (Test-Path -LiteralPath $Backup)) { Move-Item -LiteralPath $Backup -Destination $Destination }; throw }
  $Expanded = Join-Path $TempDir 'skill'; Expand-Archive -LiteralPath $Skill -DestinationPath $Expanded; $Source = Join-Path $Expanded 'guiho-s-runx'; if (-not (Test-Path -LiteralPath (Join-Path $Source 'SKILL.md'))) { throw 'Skill archive is missing guiho-s-runx/SKILL.md' }
  foreach ($Root in @((Join-Path $HOME '.agents\skills'),(Join-Path $HOME '.claude\skills'))) { New-Item -ItemType Directory -Force -Path $Root | Out-Null; $Target=Join-Path $Root 'guiho-s-runx'; if(Test-Path -LiteralPath $Target){Remove-Item -LiteralPath $Target -Recurse -Force};Copy-Item -LiteralPath $Source -Destination $Target -Recurse;Write-Host "[OK] Installed agent skill: $Target" }
  if ((Test-Path -LiteralPath 'AGENTS.md') -or (Test-Path -LiteralPath 'CLAUDE.md')) { & $Destination agent instruction update }
  if ($env:RUNX_SKIP_PATH_UPDATE -ne '1') {
    $WindowsPath = Add-WindowsPath -InstallDirectory $InstallDir
    if ($WindowsPath.UserPathChanged) { Write-Host "[OK] Added to the Windows user Path: $($WindowsPath.Directory)" }
    else { Write-Host "[OK] Already in the Windows user Path: $($WindowsPath.Directory)" }
    if ($WindowsPath.ProcessPathChanged) { Write-Host '[OK] Updated Path for the current PowerShell process.' }

    $GitBashPath = Add-GitBashPath -InstallDirectory $InstallDir
    if ($GitBashPath.Changed) { Write-Host "[OK] Added to Git Bash startup: $($GitBashPath.ProfilePath)" }
    else { Write-Host "[OK] Git Bash startup already contains: $($GitBashPath.ExportLine)" }
    Write-Host 'Existing Git Bash session: source ~/.bashrc'
  }
  Write-Host "[OK] Installed and verified RunX $ActualVersion at $Destination"
} finally { if (Test-Path -LiteralPath $TempDir) { Remove-Item -LiteralPath $TempDir -Recurse -Force } }
