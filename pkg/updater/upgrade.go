package updater

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/CGuiho/runx/pkg/update"
)

const (
	maxBinaryBytes    = 256 << 20
	maxChecksumsBytes = 1 << 20
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
		currentVer = "dev"
	}

	goos := opts.GOOS
	if goos == "" {
		goos = runtime.GOOS
	}
	goarch := opts.GOARCH
	if goarch == "" {
		goarch = runtime.GOARCH
	}

	platform, err := update.ResolveBuildTarget(opts.BuildTarget, goos, goarch)
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

	client := opts.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Minute}
	}
	catalog, err := update.FetchReleaseCatalog(opts.APIURL, platform, currentVer, client)
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
	reqVer := strings.TrimPrefix(opts.RequestedVersion, "@guiho/runx/v")
	reqVer = strings.TrimPrefix(reqVer, "@guiho/runx@")
	reqVer = strings.TrimPrefix(reqVer, "runx/v")
	reqVer = strings.TrimPrefix(reqVer, "runx@")
	reqVer = strings.TrimPrefix(reqVer, "v")

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

	if (targetEntry.CompatibleAsset == nil || targetEntry.ChecksumsAsset == nil) && targetVersion != currentVer {
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
				Message: fmt.Sprintf("RunX %s has no compatible %s asset and checksum manifest", targetVersion, platform.Target),
			},
		}, nil
	}

	plan := &UpgradePlan{
		CurrentVersion: currentVer,
		TargetVersion:  targetVersion,
		OS:             platform.OS,
		Arch:           platform.Arch,
		ExecutablePath: execPath,
		BuildTarget:    platform.Target,
	}
	if targetEntry.CompatibleAsset != nil {
		plan.AssetName = targetEntry.CompatibleAsset.Name
		plan.AssetURL = targetEntry.CompatibleAsset.URL
	}
	if targetEntry.ChecksumsAsset != nil {
		plan.ChecksumsURL = targetEntry.ChecksumsAsset.URL
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
		resp, dErr := client.Get(plan.AssetURL)
		if dErr != nil {
			err = dErr
		} else {
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				err = fmt.Errorf("HTTP %d", resp.StatusCode)
			} else {
				binaryData, err = readBounded(resp.Body, maxBinaryBytes, "release executable")
			}
		}
	}
	if err == nil && len(binaryData) > maxBinaryBytes {
		err = fmt.Errorf("release executable exceeds %d bytes", maxBinaryBytes)
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
	var checksumsData []byte
	if opts.DownloadFunc != nil {
		checksumsData, err = opts.DownloadFunc(plan.ChecksumsURL)
	} else {
		response, downloadErr := client.Get(plan.ChecksumsURL)
		if downloadErr != nil {
			err = downloadErr
		} else {
			defer response.Body.Close()
			if response.StatusCode != http.StatusOK {
				err = fmt.Errorf("HTTP %d", response.StatusCode)
			} else {
				checksumsData, err = readBounded(response.Body, maxChecksumsBytes, "checksum manifest")
			}
		}
	}
	if err == nil && len(checksumsData) > maxChecksumsBytes {
		err = fmt.Errorf("checksum manifest exceeds %d bytes", maxChecksumsBytes)
	}
	if err != nil {
		emit("validate", "failed", err.Error())
		return failedEnvelope(plan, events, recovery, "checksum_download_failed", "validate", err.Error()), nil
	}
	expected, err := checksumForAsset(checksumsData, plan.AssetName)
	if err != nil {
		emit("validate", "failed", err.Error())
		return failedEnvelope(plan, events, recovery, "checksum_missing", "validate", err.Error()), nil
	}
	hash := sha256.Sum256(binaryData)
	actual := hex.EncodeToString(hash[:])
	plan.ExpectedSHA256 = expected
	if !strings.EqualFold(actual, expected) {
		message := fmt.Sprintf("checksum mismatch for %s", plan.AssetName)
		emit("validate", "failed", message)
		return failedEnvelope(plan, events, recovery, "checksum_mismatch", "validate", message), nil
	}
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
		opts.MaintenanceCWD,
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

	outcome := "upgraded"
	if res.CleanupDeferred && platform.OS == "windows" {
		outcome = "scheduled"
	} else if opts.Spawn != nil && opts.MaintenanceCWD != "" {
		_ = opts.Spawn(execPath, "__maintenance-worker", "--cwd", opts.MaintenanceCWD)
	}
	return &UpgradeEnvelope{
		SchemaVersion: 1,
		Command:       "runx upgrade",
		Outcome:       outcome,
		Plan:          plan,
		Events:        events,
		Result:        res,
		Recovery:      recovery,
	}, nil
}

func readBounded(reader io.Reader, maximum int64, label string) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(reader, maximum+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > maximum {
		return nil, fmt.Errorf("%s exceeds %d bytes", label, maximum)
	}
	return data, nil
}

func checksumForAsset(data []byte, asset string) (string, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 2 && strings.TrimPrefix(fields[1], "*") == asset {
			if _, err := hex.DecodeString(fields[0]); err != nil || len(fields[0]) != 64 {
				return "", fmt.Errorf("invalid SHA-256 entry for %s", asset)
			}
			return strings.ToLower(fields[0]), nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("checksum entry missing for %s", asset)
}

func failedEnvelope(plan *UpgradePlan, events []UpgradeEvent, recovery RecoveryInstructions, code, phase, message string) *UpgradeEnvelope {
	return &UpgradeEnvelope{SchemaVersion: 1, Command: "runx upgrade", Outcome: "failed", Plan: plan, Events: events, Recovery: recovery, Error: &UpgradeError{Code: code, Phase: phase, Message: message}}
}
