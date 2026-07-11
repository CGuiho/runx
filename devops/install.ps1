$ErrorActionPreference = 'Stop'

$architecture = if ([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString().ToLowerInvariant() -eq 'arm64') { 'arm64' } else { 'x64' }
$assetName = "runx-windows-$architecture.exe"
$release = Invoke-RestMethod -Uri 'https://api.github.com/repos/CGuiho/runx/releases/latest' -Headers @{ Accept = 'application/vnd.github+json' }
$asset = $release.assets | Where-Object { $_.name -eq $assetName } | Select-Object -First 1
if ($null -eq $asset) { throw "Release asset $assetName was not found." }

$directory = Join-Path $HOME '.local\bin'
$destination = Join-Path $directory 'runx.exe'
New-Item -ItemType Directory -Force -Path $directory | Out-Null
Invoke-WebRequest -Uri $asset.browser_download_url -OutFile $destination

Write-Output "Installed RunX $($release.tag_name) to $destination"
if (($env:Path -split ';') -notcontains $directory) {
  Write-Output "Add $directory to your user PATH, then run: runx --help"
}
