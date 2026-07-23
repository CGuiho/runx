package updater

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/CGuiho/runx/pkg/update"
)

func CreateRecoveryInstructions(targetVersion, osName, targetSource string) RecoveryInstructions {
	if osName == "windows" {
		return RecoveryInstructions{
			TargetVersion:      targetVersion,
			TargetSource:       targetSource,
			InstallCommand:     fmt.Sprintf(`powershell.exe -NoProfile -ExecutionPolicy Bypass -Command '& ([scriptblock]::Create((Invoke-RestMethod "https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.ps1"))) -Version "%s"'`, targetVersion),
			StopProcessCommand: `powershell.exe -NoProfile -Command "Get-Process runx -ErrorAction SilentlyContinue | Stop-Process -Force"`,
		}
	}
	return RecoveryInstructions{
		TargetVersion:      targetVersion,
		TargetSource:       targetSource,
		InstallCommand:     fmt.Sprintf(`curl -fsSL https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.sh | bash -s -- --version '%s'`, targetVersion),
		StopProcessCommand: `pkill -x runx`,
	}
}

func UpgradeSelf(opts UpgradeOptions) (*UpgradeEnvelope, error) {
	currentVer := opts.CurrentVersion
	if currentVer == "" {
		currentVer = "0.8.0-dev"
	}

	goos := opts.GOOS
	if goos == "" {
		goos = runtime.GOOS
	}
	goarch := opts.GOARCH
	if goarch == "" {
		goarch = runtime.GOARCH
	}

	platform, err := update.ResolveUpgradePlatform(goos, goarch)
	if err != nil {
		rec := CreateRecoveryInstructions(currentVer, goos, "fallback-current")
		return &UpgradeEnvelope{
			SchemaVersion: 1,
			Command:       "runx upgrade",
			Outcome:       "failed",
			Recovery:      rec,
			Error: &UpgradeError{
				Code:    "platform_unsupported",
				Phase:   "plan",
				Message: err.Error(),
			},
		}, nil
	}

	execPath := opts.ExecutablePath
	if execPath == "" {
		execPath, _ = os.Executable()
	}

	recovery := CreateRecoveryInstructions(currentVer, platform.OS, "fallback-current")
	var events []UpgradeEvent

	emit := func(phase, status, message string) {
		ev := UpgradeEvent{
			Sequence: len(events) + 1,
			Phase:    phase,
			Status:   status,
			Message:  message,
		}
		events = append(events, ev)
		if opts.OnEvent != nil {
			opts.OnEvent(ev)
		}
	}

	emit("plan", "started", "")

	catalog, err := update.FetchReleaseCatalog(opts.APIURL, platform, currentVer, nil)
	if err != nil {
		emit("plan", "failed", err.Error())
		return &UpgradeEnvelope{
			SchemaVersion: 1,
			Command:       "runx upgrade",
			Outcome:       "failed",
			Events:        events,
			Recovery:      recovery,
			Error: &UpgradeError{
				Code:    "release_lookup_failed",
				Phase:   "plan",
				Message: err.Error(),
			},
		}, nil
	}

	var targetEntry *update.ReleaseCatalogEntry
	reqVer := strings.TrimPrefix(opts.RequestedVersion, "v")
	reqVer = strings.TrimPrefix(reqVer, "@guiho/runx@")

	if reqVer != "" {
		for i := range catalog.Releases {
			if catalog.Releases[i].Version == reqVer {
				targetEntry = &catalog.Releases[i]
				break
			}
		}
	} else {
		for i := range catalog.Releases {
			if catalog.Releases[i].LatestStable {
				targetEntry = &catalog.Releases[i]
				break
			}
		}
	}

	if targetEntry == nil {
		emit("plan", "failed", "No compatible release target found")
		return &UpgradeEnvelope{
			SchemaVersion: 1,
			Command:       "runx upgrade",
			Outcome:       "failed",
			Events:        events,
			Recovery:      recovery,
			Error: &UpgradeError{
				Code:    "release_lookup_failed",
				Phase:   "plan",
				Message: "No compatible release target found",
			},
		}, nil
	}

	targetVersion := targetEntry.Version
	recovery = CreateRecoveryInstructions(targetVersion, platform.OS, "resolved")

	if targetEntry.CompatibleAsset == nil && targetVersion != currentVer {
		emit("plan", "failed", "No compatible asset found")
		return &UpgradeEnvelope{
			SchemaVersion: 1,
			Command:       "runx upgrade",
			Outcome:       "failed",
			Events:        events,
			Recovery:      recovery,
			Error: &UpgradeError{
				Code:    "no_compatible_asset",
				Phase:   "plan",
				Message: fmt.Sprintf("RunX %s has no compatible %s %s asset", targetVersion, platform.OS, platform.Arch),
			},
		}, nil
	}

	plan := &UpgradePlan{
		CurrentVersion: currentVer,
		TargetVersion:  targetVersion,
		OS:             platform.OS,
		Arch:           platform.Arch,
		ExecutablePath: execPath,
	}
	if targetEntry.CompatibleAsset != nil {
		plan.AssetName = targetEntry.CompatibleAsset.Name
		plan.AssetURL = targetEntry.CompatibleAsset.URL
	}

	emit("plan", "succeeded", "")
	if opts.OnPlan != nil {
		opts.OnPlan(*plan)
	}

	if update.CompareVersions(targetVersion, currentVer) <= 0 && reqVer == "" {
		return &UpgradeEnvelope{
			SchemaVersion: 1,
			Command:       "runx upgrade",
			Outcome:       "up-to-date",
			Plan:          plan,
			Events:        events,
			Result: &UpgradeResult{
				InstalledVersion: currentVer,
				CleanupDeferred:  false,
			},
			Recovery: recovery,
		}, nil
	}

	if opts.DryRun {
		return &UpgradeEnvelope{
			SchemaVersion: 1,
			Command:       "runx upgrade",
			Outcome:       "dry-run",
			Plan:          plan,
			Events:        events,
			Recovery:      recovery,
		}, nil
	}

	// Download phase
	emit("download", "started", "")
	var binaryData []byte

	if opts.DownloadFunc != nil {
		binaryData, err = opts.DownloadFunc(plan.AssetURL)
	} else {
		resp, dErr := http.Get(plan.AssetURL)
		if dErr != nil {
			err = dErr
		} else {
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				err = fmt.Errorf("HTTP %d", resp.StatusCode)
			} else {
				binaryData, err = io.ReadAll(resp.Body)
			}
		}
	}

	if err != nil {
		emit("download", "failed", err.Error())
		return &UpgradeEnvelope{
			SchemaVersion: 1,
			Command:       "runx upgrade",
			Outcome:       "failed",
			Plan:          plan,
			Events:        events,
			Recovery:      recovery,
			Error: &UpgradeError{
				Code:    "download_failed",
				Phase:   "download",
				Message: err.Error(),
			},
		}, nil
	}
	emit("download", "succeeded", "")

	// Validate phase
	emit("validate", "started", "")
	if valErr := ValidateNativeBinary(binaryData, platform.OS); valErr != nil {
		emit("validate", "failed", valErr.Error())
		return &UpgradeEnvelope{
			SchemaVersion: 1,
			Command:       "runx upgrade",
			Outcome:       "failed",
			Plan:          plan,
			Events:        events,
			Recovery:      recovery,
			Error: &UpgradeError{
				Code:    "download_invalid",
				Phase:   "validate",
				Message: valErr.Error(),
			},
		}, nil
	}
	emit("validate", "succeeded", "")

	// Replace & Verify phase
	emit("replace", "started", "")

	fileOps := opts.FileOps
	if fileOps == nil {
		fileOps = SystemFileOps{}
	}

	res, errCode, repErr := PerformReplacementAndRollback(
		execPath,
		binaryData,
		targetVersion,
		platform.OS,
		fileOps,
		opts.VerifyFunc,
	)

	if repErr != nil {
		if errCode == "rollback_failed" {
			emit("replace", "failed", repErr.Error())
			return &UpgradeEnvelope{
				SchemaVersion: 1,
				Command:       "runx upgrade",
				Outcome:       "failed",
				Plan:          plan,
				Events:        events,
				Recovery:      recovery,
				Error: &UpgradeError{
					Code:    "rollback_failed",
					Phase:   "replace",
					Message: repErr.Error(),
				},
			}, nil
		}

		emit("replace", "failed", repErr.Error())
		return &UpgradeEnvelope{
			SchemaVersion: 1,
			Command:       "runx upgrade",
			Outcome:       "rolled-back",
			Plan:          plan,
			Events:        events,
			Result: &UpgradeResult{
				InstalledVersion: currentVer,
				CleanupDeferred:  false,
			},
			Recovery: recovery,
			Error: &UpgradeError{
				Code:    errCode,
				Phase:   "replace",
				Message: repErr.Error(),
			},
		}, nil
	}

	emit("replace", "succeeded", "")
	emit("verify", "started", "")
	emit("verify", "succeeded", "")
	emit("cleanup", "started", "")
	emit("cleanup", "succeeded", "")

	return &UpgradeEnvelope{
		SchemaVersion: 1,
		Command:       "runx upgrade",
		Outcome:       "upgraded",
		Plan:          plan,
		Events:        events,
		Result:        res,
		Recovery:      recovery,
	}, nil
}
