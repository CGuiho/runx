param(
  [string]$Version = $(if ($env:RUNX_VERSION) { $env:RUNX_VERSION } else { 'latest' }),
  [string]$InstallDir = $(if ($env:RUNX_INSTALL_DIR) { $env:RUNX_INSTALL_DIR } else { Join-Path $HOME '.local\bin' }),
  [switch]$Help
)

$ErrorActionPreference = 'Stop'
if ($Help) { Write-Output 'Usage: install.ps1 [-Version VERSION] [-InstallDir DIRECTORY]'; return }
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
  if ($env:RUNX_SKIP_PATH_UPDATE -ne '1') { $UserPath=[Environment]::GetEnvironmentVariable('Path','User');if(-not (($UserPath -split ';') -contains $InstallDir)){[Environment]::SetEnvironmentVariable('Path',(($InstallDir,$UserPath -join ';').TrimEnd(';')),'User')} }
  Write-Host "[OK] Installed and verified RunX $ActualVersion at $Destination"
} finally { if (Test-Path -LiteralPath $TempDir) { Remove-Item -LiteralPath $TempDir -Recurse -Force } }
