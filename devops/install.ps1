# GUIHO RunX installer (Convention 0001).
# Installs the stable launcher into $HOME\.guiho\bin\, immutable payloads into
# $HOME\.guiho\runx\versions\<version>\, activates via an atomic current.json
# pointer, verifies every artifact against checksums.txt and artifacts.json,
# and rolls back completely on failure.
param(
  [string]$Version = '',
  [string]$Channel = '',
  [switch]$Help
)

$ErrorActionPreference = 'Stop'

if ($Help) { Write-Output 'Usage: install.ps1 [-Version VERSION] [-Channel CHANNEL]'; return }
if ($Version -and $Channel) { throw '--Version and -Channel are mutually exclusive.' }

$Repo = if ($env:RUNX_REPO) { $env:RUNX_REPO } else { 'CGuiho/runx' }
$Platform = switch ($env:PROCESSOR_ARCHITECTURE.ToUpperInvariant()) {
  'AMD64' { 'windows-amd64' }
  'ARM64' { 'windows-arm64' }
  default { throw "Unsupported Windows architecture: $env:PROCESSOR_ARCHITECTURE" }
}
$PayloadAsset = "runx-payload-$Platform.exe"
$LauncherAsset = "runx-launcher-$Platform.exe"

function Normalize-Version([string]$Value) {
  return ($Value -replace '^@guiho/runx/v', '' -replace '^@guiho/runx@', '' -replace '^runx/v', '' -replace '^runx@', '' -replace '^v', '')
}

function Test-SemVer([string]$Value) {
  return $Value -match '^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$'
}

function Get-Channel([string]$SemVer) {
  if ($SemVer -match '-') { return (($SemVer -split '-')[1] -split '\.')[0] }
  return 'stable'
}

$Headers = @{ Accept = 'application/vnd.github+json' }

function Resolve-Selection {
  if ($Version) {
    $Normalized = Normalize-Version $Version
    if (-not (Test-SemVer $Normalized)) { throw "invalid RunX version: $Version" }
    return $Normalized
  }
  for ($Page = 1; $Page -le 10; $Page++) {
    $Response = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases?per_page=100&page=$Page" -Headers $Headers
    foreach ($Release in $Response) {
      $Normalized = Normalize-Version $Release.tag_name
      if (-not (Test-SemVer $Normalized)) { continue }
      $CandidateChannel = if ($Release.prerelease) { Get-Channel $Normalized } else { 'stable' }
      if (-not $Channel -or $Channel -eq 'stable') {
        if ($CandidateChannel -eq 'stable') { return $Normalized }
      } elseif ($CandidateChannel -eq $Channel) {
        return $Normalized
      }
    }
    if ($Response.Count -lt 100) { break }
  }
  if ($Channel) { throw "no published release found for channel: $Channel" }
  throw 'no published stable release found'
}

$TargetVersion = Resolve-Selection
# Transition: prefer new tag runx/v* if it exists, fallback to legacy @guiho/runx/v*.
$EncodedTagNew = [Uri]::EscapeDataString("runx/v$TargetVersion")
$EncodedTagOld = [Uri]::EscapeDataString("@guiho/runx/v$TargetVersion")
$ProbeNew = "https://github.com/$Repo/releases/download/$EncodedTagNew/checksums.txt"
try { Invoke-WebRequest -Uri $ProbeNew -Method Head -UseBasicParsing -TimeoutSec 10 | Out-Null; $EncodedTag = $EncodedTagNew } catch { $EncodedTag = $EncodedTagOld }
$AssetBase = "https://github.com/$Repo/releases/download/$EncodedTag"

$GuihoRoot = Join-Path $HOME '.guiho'
$CliDir = Join-Path $GuihoRoot 'runx'
$BinDir = Join-Path $GuihoRoot 'bin'
$VersionsDir = Join-Path $CliDir 'versions'
$ResourcesDir = Join-Path $CliDir 'resources'
$TempRoot = Join-Path $GuihoRoot '.temp'

New-Item -ItemType Directory -Force -Path $TempRoot | Out-Null
$Staging = Join-Path $TempRoot ("runx-install-" + [Guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path $Staging | Out-Null

Write-Host 'Initiating GUIHO CLI Upgrade / Installation Sequence...'
Write-Host "Target Version: $TargetVersion"
Write-Host "Platform:       $Platform"
Write-Host "Payload Asset:  $PayloadAsset"
Write-Host "Source URL:     $AssetBase"
Write-Host "Staging:        $Staging"

function Download([string]$Name) {
  Invoke-WebRequest -Uri "$AssetBase/$Name" -OutFile (Join-Path $Staging $Name) -UseBasicParsing
}
foreach ($Asset in @($PayloadAsset, $LauncherAsset, 'checksums.txt', 'artifacts.json', 'guiho-s-runx.zip', 'guiho-i-runx.md', 'guiho-p-runx.md', 'guiho-p-runx-uninstall.md', 'runx.schema.json', 'runx.global.schema.json')) {
  Download $Asset
}

$ChecksumLines = Get-Content -LiteralPath (Join-Path $Staging 'checksums.txt')
function Verify-Checksum([string]$Name) {
  $Line = $ChecksumLines | Where-Object { $_ -match "\s+$([Regex]::Escape($Name))$" } | Select-Object -First 1
  if (-not $Line) { throw "checksum entry missing for $Name" }
  $Expected = ($Line -split '\s+')[0]
  $Actual = (Get-FileHash -Algorithm SHA256 -LiteralPath (Join-Path $Staging $Name)).Hash.ToLowerInvariant()
  if ($Expected.ToLowerInvariant() -ne $Actual) { throw "checksum verification failed for $Name" }
}
$Manifest = Get-Content -Raw -LiteralPath (Join-Path $Staging 'artifacts.json') | ConvertFrom-Json
function Verify-ManifestDigest([string]$Name) {
  $Entry = $Manifest.artifacts | Where-Object { $_.file -eq $Name } | Select-Object -First 1
  if ($Entry) {
    $Actual = (Get-FileHash -Algorithm SHA256 -LiteralPath (Join-Path $Staging $Name)).Hash.ToLowerInvariant()
    if ($Entry.sha256.ToLowerInvariant() -ne $Actual) { throw "artifacts.json digest mismatch for $Name" }
  }
}
foreach ($Asset in @($PayloadAsset, $LauncherAsset, 'guiho-s-runx.zip', 'guiho-i-runx.md', 'guiho-p-runx.md', 'guiho-p-runx-uninstall.md', 'runx.schema.json', 'runx.global.schema.json', 'artifacts.json')) {
  Verify-Checksum $Asset
  Verify-ManifestDigest $Asset
}
Write-Host '[OK] SHA-256 verification complete.'

# Self-test the staged payload BEFORE activation.
$StagedPayload = Join-Path $Staging $PayloadAsset
$PreviousUpdate = $env:RUNX_DISABLE_UPDATE_WORKER
$PreviousMaintenance = $env:RUNX_DISABLE_AGENT_MAINTENANCE_WORKER
$env:RUNX_DISABLE_UPDATE_WORKER = '1'; $env:RUNX_DISABLE_AGENT_MAINTENANCE_WORKER = '1'
try {
  $StagedVersion = (& $StagedPayload --version).Trim()
  if ($StagedVersion -ne $TargetVersion) { throw "staged payload reports $StagedVersion, want $TargetVersion" }
  & $StagedPayload __self-test *> $null
  if ($LASTEXITCODE -ne 0) { throw 'staged payload failed its installation self-test' }
} finally {
  $env:RUNX_DISABLE_UPDATE_WORKER = $PreviousUpdate; $env:RUNX_DISABLE_AGENT_MAINTENANCE_WORKER = $PreviousMaintenance
}

$PointerPath = Join-Path $CliDir 'current.json'
$HadPointer = Test-Path -LiteralPath $PointerPath
if ($HadPointer) { Copy-Item -LiteralPath $PointerPath -Destination (Join-Path $Staging 'current.json.previous') }
function Restore-Pointer {
  if ($HadPointer) { Copy-Item -Force -LiteralPath (Join-Path $Staging 'current.json.previous') -Destination $PointerPath }
  else { Remove-Item -Force -ErrorAction SilentlyContinue -LiteralPath $PointerPath }
}

$DestVersionDir = Join-Path $VersionsDir $TargetVersion
if (Test-Path -LiteralPath $DestVersionDir) { Remove-Item -Recurse -Force -LiteralPath $DestVersionDir }
New-Item -ItemType Directory -Force -Path $DestVersionDir | Out-Null
Copy-Item -LiteralPath $StagedPayload -Destination (Join-Path $DestVersionDir 'runx-payload.exe')
Copy-Item -LiteralPath (Join-Path $Staging 'artifacts.json') -Destination (Join-Path $DestVersionDir 'release-artifacts.json')

New-Item -ItemType Directory -Force -Path (Join-Path $ResourcesDir 'skills'), (Join-Path $ResourcesDir 'instruction'), (Join-Path $ResourcesDir 'prompts'), (Join-Path $ResourcesDir 'schemas') | Out-Null
Expand-Archive -LiteralPath (Join-Path $Staging 'guiho-s-runx.zip') -DestinationPath (Join-Path $ResourcesDir 'skills') -Force
if (-not (Test-Path -LiteralPath (Join-Path $ResourcesDir 'skills\guiho-s-runx\SKILL.md'))) { Restore-Pointer; Remove-Item -Recurse -Force -LiteralPath $DestVersionDir; throw 'skill archive is missing guiho-s-runx/SKILL.md' }
Copy-Item -LiteralPath (Join-Path $Staging 'guiho-i-runx.md') -Destination (Join-Path $ResourcesDir 'instruction\guiho-i-runx.md') -Force
Copy-Item -LiteralPath (Join-Path $Staging 'guiho-p-runx.md') -Destination (Join-Path $ResourcesDir 'prompts\guiho-p-runx.md') -Force
Copy-Item -LiteralPath (Join-Path $Staging 'guiho-p-runx-uninstall.md') -Destination (Join-Path $ResourcesDir 'prompts\guiho-p-runx-uninstall.md') -Force
Copy-Item -LiteralPath (Join-Path $Staging 'runx.schema.json') -Destination (Join-Path $ResourcesDir 'schemas\runx.schema.json') -Force
Copy-Item -LiteralPath (Join-Path $Staging 'runx.global.schema.json') -Destination (Join-Path $ResourcesDir 'schemas\runx.global.schema.json') -Force

foreach ($Root in @((Join-Path $HOME '.agents\skills'), (Join-Path $HOME '.claude\skills'))) {
  New-Item -ItemType Directory -Force -Path $Root | Out-Null
  $SkillTarget = Join-Path $Root 'guiho-s-runx'
  if (Test-Path -LiteralPath $SkillTarget) { Remove-Item -LiteralPath $SkillTarget -Recurse -Force }
  Copy-Item -LiteralPath (Join-Path $ResourcesDir 'skills\guiho-s-runx') -Destination $SkillTarget -Recurse
  Write-Host "[OK] Installed agent skill: $SkillTarget"
}
if ((Test-Path -LiteralPath 'AGENTS.md') -or (Test-Path -LiteralPath 'CLAUDE.md')) {
  & (Join-Path $DestVersionDir 'runx-payload.exe') agent instruction update
}

$PreviousActive = ''
if ($HadPointer) {
  $OldPointer = Get-Content -Raw -LiteralPath (Join-Path $Staging 'current.json.previous') | ConvertFrom-Json
  $PreviousActive = $OldPointer.active
}
if ($PreviousActive -and $PreviousActive -ne $TargetVersion) {
  @{ protocol = 1; active = $TargetVersion; previous = $PreviousActive } |
    ConvertTo-Json -Compress | Set-Content -NoNewline -Encoding utf8 (Join-Path $Staging 'current.json.new')
} else {
  @{ protocol = 1; active = $TargetVersion } |
    ConvertTo-Json -Compress | Set-Content -NoNewline -Encoding utf8 (Join-Path $Staging 'current.json.new')
}

$LauncherPath = Join-Path $BinDir 'runx.exe'
$BackupLauncher = Join-Path $Staging 'launcher.previous'
$HadLauncher = Test-Path -LiteralPath $LauncherPath
if ($HadLauncher) { Copy-Item -LiteralPath $LauncherPath -Destination $BackupLauncher }

function Commit-Activation {
  New-Item -ItemType Directory -Force -Path $BinDir | Out-Null
  Move-Item -Force -LiteralPath (Join-Path $Staging 'current.json.new') -Destination $PointerPath
  Copy-Item -LiteralPath (Join-Path $Staging $LauncherAsset) -Destination "$LauncherPath.tmp" -Force
  Move-Item -Force -LiteralPath "$LauncherPath.tmp" -Destination $LauncherPath
}
function Restore-Installation {
  Restore-Pointer
  if ($HadLauncher) { Copy-Item -Force -LiteralPath $BackupLauncher -Destination $LauncherPath }
  Remove-Item -Recurse -Force -ErrorAction SilentlyContinue -LiteralPath $DestVersionDir
}

try { Commit-Activation } catch { Restore-Installation; throw 'activation failed and was rolled back' }

$env:RUNX_DISABLE_UPDATE_WORKER = '1'; $env:RUNX_DISABLE_AGENT_MAINTENANCE_WORKER = '1'
try { $ActualVersion = (& $LauncherPath --version).Trim() } catch { $ActualVersion = '' } finally {
  $env:RUNX_DISABLE_UPDATE_WORKER = $PreviousUpdate; $env:RUNX_DISABLE_AGENT_MAINTENANCE_WORKER = $PreviousMaintenance
}
if ($ActualVersion -ne $TargetVersion) { Restore-Installation; throw "activated launcher reports '$ActualVersion', want $TargetVersion" }

# User-level PATH update: idempotent, no admin rights.
function Get-ComparablePath([string]$Value) {
  try { return [IO.Path]::GetFullPath(([Environment]::ExpandEnvironmentVariables($Value.Trim().Trim('"')))).TrimEnd('\', '/') }
  catch { return $Value.Trim().Trim('"').TrimEnd('\', '/') }
}
$UserPath = [Environment]::GetEnvironmentVariable('Path', 'User')
$ComparableBin = Get-ComparablePath $BinDir
$AlreadyInPath = $false
foreach ($Entry in ($UserPath -split ';')) {
  if ($Entry -and (Get-ComparablePath $Entry) -ieq $ComparableBin) { $AlreadyInPath = $true; break }
}
if (-not $AlreadyInPath) {
  $Updated = if ([string]::IsNullOrWhiteSpace($UserPath)) { $ComparableBin } else { "$ComparableBin;$UserPath" }
  [Environment]::SetEnvironmentVariable('Path', $Updated, 'User')
  if ($env:Path -notlike "*$ComparableBin*") { $env:Path = "$ComparableBin;$env:Path" }
  Write-Host "[OK] Added to the Windows user Path: $ComparableBin"
} else {
  Write-Host "[OK] Already in the Windows user Path: $ComparableBin"
}

Remove-Item -Recurse -Force -ErrorAction SilentlyContinue -LiteralPath $Staging
Write-Host "[OK] Installed and verified RunX $ActualVersion"
Write-Host "Launcher: $LauncherPath"
Write-Host "Active payload: $(Join-Path $DestVersionDir 'runx-payload.exe')"
Write-Host "RunX home: $CliDir"
