//go:build windows

package updater

import (
	"encoding/base64"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"unicode/utf16"
)

func stageWindowsReplacement(pid int, destination, candidate, expectedVersion, maintenanceCWD string) error {
	script := `$ErrorActionPreference='Stop'; $pidToWait=[int]$env:RUNX_PARENT_PID; for($i=0;$i -lt 600;$i++){if(-not(Get-Process -Id $pidToWait -ErrorAction SilentlyContinue)){break};Start-Sleep -Milliseconds 100}; if(Get-Process -Id $pidToWait -ErrorAction SilentlyContinue){throw 'RunX process did not exit'}; $backup="$env:RUNX_DESTINATION.old-$pidToWait"; try { Move-Item -LiteralPath $env:RUNX_DESTINATION -Destination $backup -Force; Move-Item -LiteralPath $env:RUNX_CANDIDATE -Destination $env:RUNX_DESTINATION -Force; $psi=[Diagnostics.ProcessStartInfo]::new();$psi.FileName=$env:RUNX_DESTINATION;$psi.Arguments='--version';$psi.UseShellExecute=$false;$psi.CreateNoWindow=$true;$psi.RedirectStandardOutput=$true;$psi.RedirectStandardError=$true;$psi.EnvironmentVariables['RUNX_DISABLE_UPDATE_WORKER']='1';$psi.EnvironmentVariables['RUNX_DISABLE_AGENT_MAINTENANCE_WORKER']='1';$p=[Diagnostics.Process]::new();$p.StartInfo=$psi;[void]$p.Start();if(-not $p.WaitForExit(10000)){$p.Kill();throw 'replacement version check timed out'};$actual=$p.StandardOutput.ReadToEnd().Trim();if($p.ExitCode -ne 0 -or $actual -ne $env:RUNX_EXPECTED_VERSION){throw "replacement reported $actual"};Remove-Item -LiteralPath $backup -Force -ErrorAction SilentlyContinue;if($env:RUNX_MAINTENANCE_CWD){$worker=[Diagnostics.ProcessStartInfo]::new();$worker.FileName=$env:RUNX_DESTINATION;$worker.Arguments='__maintenance-worker';$worker.UseShellExecute=$false;$worker.CreateNoWindow=$true;$worker.EnvironmentVariables['RUNX_MAINTENANCE_CWD']=$env:RUNX_MAINTENANCE_CWD;$worker.EnvironmentVariables['RUNX_DISABLE_UPDATE_WORKER']='1';$worker.EnvironmentVariables['RUNX_DISABLE_AGENT_MAINTENANCE_WORKER']='1';[void][Diagnostics.Process]::Start($worker)} } catch { if(Test-Path -LiteralPath $env:RUNX_DESTINATION){Remove-Item -LiteralPath $env:RUNX_DESTINATION -Force -ErrorAction SilentlyContinue};if(Test-Path -LiteralPath $backup){Move-Item -LiteralPath $backup -Destination $env:RUNX_DESTINATION -Force};Set-Content -LiteralPath "$env:RUNX_DESTINATION.upgrade-error.log" -Value $_.Exception.Message -Encoding UTF8;exit 1 }`
	runes := utf16.Encode([]rune(script))
	encodedBytes := make([]byte, len(runes)*2)
	for index, value := range runes {
		encodedBytes[index*2] = byte(value)
		encodedBytes[index*2+1] = byte(value >> 8)
	}
	encoded := base64.StdEncoding.EncodeToString(encodedBytes)
	command := exec.Command("powershell.exe", "-NoLogo", "-NoProfile", "-NonInteractive", "-EncodedCommand", encoded)
	command.Env = append(os.Environ(), "RUNX_PARENT_PID="+strconv.Itoa(pid), "RUNX_DESTINATION="+destination, "RUNX_CANDIDATE="+candidate, "RUNX_EXPECTED_VERSION="+expectedVersion, "RUNX_MAINTENANCE_CWD="+maintenanceCWD)
	command.Stdin, command.Stdout, command.Stderr = nil, nil, nil
	command.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x00000008 | 0x08000000}
	return command.Start()
}
