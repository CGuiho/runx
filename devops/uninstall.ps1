# GUIHO RunX uninstaller (Convention 0001).
# Shares the uninstallation contract with `runx uninstall` and
# devops/uninstall.sh. By default removes everything RunX owns; supports
# -PreserveConfig, -PreserveData, -DryRun, and -Yes.
param(
  [switch]$PreserveConfig,
  [switch]$PreserveData,
  [switch]$DryRun,
  [switch]$Yes,
  [switch]$Help
)

$ErrorActionPreference = 'Stop'

if ($Help) { Write-Output 'Usage: uninstall.ps1 [-PreserveConfig] [-PreserveData] [-DryRun] [-Yes]'; return }

$GuihoRoot = Join-Path $HOME '.guiho'
$CliDir = Join-Path $GuihoRoot 'runx'
$BinDir = Join-Path $GuihoRoot 'bin'
$LauncherPath = Join-Path $BinDir 'runx.exe'
$ProjectDir = (Get-Location).Path

$Remove = New-Object System.Collections.Generic.List[string]
$Preserve = New-Object System.Collections.Generic.List[string]

# Shared infrastructure is never removed.
$Preserve.Add("$GuihoRoot (shared GUIHO home)")
if (Test-Path -LiteralPath $BinDir) { $Preserve.Add("$BinDir (shared launcher directory)") }
$Preserve.Add("$(Join-Path $GuihoRoot '.temp') (shared staging root)")
$Preserve.Add("user Path entry for $BinDir")

if (Test-Path -LiteralPath $LauncherPath) { $Remove.Add($LauncherPath) }
if (Test-Path -LiteralPath $CliDir) {
  if ($PreserveConfig -or $PreserveData) {
    $Preserve.Add("$(Join-Path $CliDir 'runx.global.yaml')")
    $Preserve.Add("$(Join-Path $CliDir 'data')")
    $Preserve.Add("$CliDir directory itself (preserved children remain)")
    $Remove.Add($(Join-Path $CliDir 'versions'))
    $Remove.Add($(Join-Path $CliDir 'resources'))
    $Remove.Add($(Join-Path $CliDir 'current.json'))
    $Remove.Add($(Join-Path $CliDir 'installed-artifacts.json'))
  } else {
    $Remove.Add($CliDir)
  }
}
foreach ($Root in @((Join-Path $HOME '.agents\skills'), (Join-Path $HOME '.claude\skills'))) {
  $SkillTarget = Join-Path $Root 'guiho-s-runx'
  if (Test-Path -LiteralPath $SkillTarget) { $Remove.Add($SkillTarget) }
}
foreach ($Marker in @('AGENTS.md', 'CLAUDE.md')) {
  if (Test-Path -LiteralPath (Join-Path $ProjectDir $Marker)) { $Remove.Add("RunX managed block in $Marker") }
}
$ProjectManifest = Join-Path $ProjectDir 'runx.yaml'
if (Test-Path -LiteralPath $ProjectManifest) {
  if ($PreserveConfig) { $Preserve.Add($ProjectManifest) } else { $Remove.Add($ProjectManifest) }
}

Write-Host 'Uninstallation plan'
Write-Host ''
Write-Host 'REMOVE:'
if ($Remove.Count -eq 0) { Write-Host '  (nothing)' } else { foreach ($Item in $Remove) { Write-Host "  $Item" } }
Write-Host ''
Write-Host 'PRESERVE:'
foreach ($Item in $Preserve) { Write-Host "  $Item" }
Write-Host ''

if ($DryRun) { Write-Host 'Dry run: nothing was changed.'; return }

if (-not $Yes) {
  if (-not ([Environment]::UserInteractive)) { throw 'noninteractive uninstall requires -Yes' }
  $Answer = Read-Host 'Proceed with uninstallation? [y/N]'
  if ($Answer -notmatch '^(y|Y|yes|YES)$') { Write-Host 'Aborted.'; return }
}

function Remove-OwnedPath([string]$Target) {
  $Full = [IO.Path]::GetFullPath($Target).TrimEnd('\', '/')
  $FullCli = [IO.Path]::GetFullPath($CliDir).TrimEnd('\', '/')
  $FullLauncher = [IO.Path]::GetFullPath($LauncherPath)
  $AllowedSkills = @(
    [IO.Path]::GetFullPath((Join-Path $HOME '.agents\skills\guiho-s-runx')),
    [IO.Path]::GetFullPath((Join-Path $HOME '.claude\skills\guiho-s-runx'))
  )
  $Owned = ($Full -ieq $FullCli) -or ($Full -like "$FullCli\*") -or ($Full -ieq $FullLauncher) -or (($AllowedSkills | Where-Object { $_ -ieq $Full }).Count -gt 0)
  if (-not $Owned) { throw "refusing to remove unowned path: $Target" }
  if (Test-Path -LiteralPath $Target) {
    # Windows cannot delete a running executable; quarantine by rename first so
    # the launcher stops resolving even when deletion must complete on reboot.
    try {
      Remove-Item -Recurse -Force -LiteralPath $Target
    } catch {
      $Quarantine = "$Target.uninstall-quarantine-" + [Guid]::NewGuid().ToString('N')
      Move-Item -Force -LiteralPath $Target -Destination $Quarantine
      Write-Host "[PENDING] Quarantined for deletion after reboot: $Quarantine"
      reg add "HKCU\Software\Microsoft\Windows\CurrentVersion\RunOnce" /v "runx-uninstall-cleanup" /d "cmd.exe /c rd /s /q `"$Quarantine`"" /f | Out-Null
    }
  }
}

if (Test-Path -LiteralPath $LauncherPath) { Remove-OwnedPath $LauncherPath }
if (Test-Path -LiteralPath $CliDir) {
  if (-not ($PreserveConfig -or $PreserveData)) {
    Remove-OwnedPath $CliDir
  } else {
    Remove-OwnedPath (Join-Path $CliDir 'versions')
    Remove-OwnedPath (Join-Path $CliDir 'resources')
    foreach ($File in @('current.json', 'installed-artifacts.json')) {
      $Path = Join-Path $CliDir $File
      if (Test-Path -LiteralPath $Path) { Remove-Item -Force -LiteralPath $Path }
    }
    Get-ChildItem -LiteralPath $CliDir -Force |
      Where-Object { $_.Name -ne 'runx.global.yaml' -and $_.Name -ne 'data' } |
      ForEach-Object { Remove-Item -Recurse -Force -LiteralPath $_.FullName }
    try { Remove-Item -LiteralPath $CliDir } catch { }
  }
}
foreach ($Root in @((Join-Path $HOME '.agents\skills'), (Join-Path $HOME '.claude\skills'))) {
  $SkillTarget = Join-Path $Root 'guiho-s-runx'
  if (Test-Path -LiteralPath $SkillTarget) { Remove-OwnedPath $SkillTarget }
}

function Remove-ManagedBlock([string]$File) {
  if (-not (Test-Path -LiteralPath $File)) { return }
  $Lines = Get-Content -LiteralPath $File
  $Inside = $false
  $Kept = foreach ($Line in $Lines) {
    if ($Line -match '^<!-- BEGIN RUNX') { $Inside = $true; continue }
    if ($Line -match '^<!-- END RUNX') { $Inside = $false; continue }
    if (-not $Inside) { $Line }
  }
  Set-Content -LiteralPath $File -Value $Kept -Encoding utf8
}
foreach ($Marker in @('AGENTS.md', 'CLAUDE.md')) { Remove-ManagedBlock (Join-Path $ProjectDir $Marker) }
if (-not $PreserveConfig -and (Test-Path -LiteralPath $ProjectManifest)) {
  Remove-Item -Force -LiteralPath $ProjectManifest
}

Write-Host '[OK] RunX uninstalled.'
