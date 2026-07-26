package cmd

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/CGuiho/runx/pkg/update"
	"github.com/CGuiho/runx/pkg/updater"
	"github.com/spf13/cobra"
)

func newUpgradeCommand(deps Dependencies, info BuildInfo) *cobra.Command {
	var requested, format string
	var dryRun bool
	command := &cobra.Command{Use: "upgrade", Short: "Inspect or upgrade a native RunX executable.", Args: cobra.NoArgs, RunE: func(command *cobra.Command, _ []string) error {
		if err := validateFormat(format); err != nil {
			return err
		}
		cwd, _ := deps.Getwd()
		envelope, err := updater.UpgradeSelf(updater.UpgradeOptions{DryRun: dryRun, CurrentVersion: info.Version, RequestedVersion: requested, BuildTarget: info.Target, HTTPClient: deps.HTTPClient, MaintenanceCWD: cwd, Spawn: deps.Spawn})
		if err != nil {
			return withExitCode(4, err)
		}
		if format == "json" {
			if err := writeJSON(command, envelope); err != nil {
				return err
			}
		} else {
			renderUpgrade(command, envelope)
		}
		if envelope.Outcome == "failed" || envelope.Outcome == "rolled-back" {
			code := 4
			if envelope.Error != nil && (envelope.Error.Phase == "replace" || envelope.Error.Phase == "verify") {
				code = 5
			}
			return withExitCode(code, fmt.Errorf("upgrade %s: %s", envelope.Outcome, envelope.Error.Message))
		}
		return nil
	}}
	command.Flags().StringVar(&requested, "version", "", "Select an exact release version.")
	command.Flags().BoolVar(&dryRun, "dry-run", false, "Plan without mutation.")
	command.Flags().StringVar(&format, "format", "text", "Select output format: text or json.")
	command.AddCommand(newUpgradeCheckCommand(deps, info), newUpgradeListCommand(deps, info))
	return command
}

func releaseCatalog(deps Dependencies, info BuildInfo) (*update.ReleaseCatalog, error) {
	platform, err := update.ResolveBuildTarget(info.Target, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return nil, err
	}
	return update.FetchReleaseCatalog("", platform, info.Version, deps.HTTPClient)
}

func newUpgradeCheckCommand(deps Dependencies, info BuildInfo) *cobra.Command {
	var format string
	command := &cobra.Command{Use: "check", Short: "Check whether a newer stable release exists.", Args: cobra.NoArgs, RunE: func(command *cobra.Command, _ []string) error {
		if err := validateFormat(format); err != nil {
			return err
		}
		catalog, err := releaseCatalog(deps, info)
		if err != nil {
			return withExitCode(4, err)
		}
		available := update.CompareVersions(catalog.LatestStableVersion, info.Version) > 0
		result := struct {
			CurrentVersion      string `json:"currentVersion"`
			LatestStableVersion string `json:"latestStableVersion"`
			UpdateAvailable     bool   `json:"updateAvailable"`
		}{info.Version, catalog.LatestStableVersion, available}
		if format == "json" {
			return writeJSON(command, result)
		}
		fmt.Fprintf(command.OutOrStdout(), "current: %s\nlatest stable: %s\nupdate available: %t\n", result.CurrentVersion, result.LatestStableVersion, result.UpdateAvailable)
		return nil
	}}
	command.Flags().StringVar(&format, "format", "text", "Select output format: text or json.")
	return command
}

func newUpgradeListCommand(deps Dependencies, info BuildInfo) *cobra.Command {
	var format string
	var page, size int
	command := &cobra.Command{Use: "list", Short: "List RunX releases newest first.", Args: cobra.NoArgs, RunE: func(command *cobra.Command, _ []string) error {
		if err := validateFormat(format); err != nil {
			return err
		}
		if page < 1 || size < 1 || size > 100 {
			return withExitCode(2, fmt.Errorf("--page must be positive and --size must be between 1 and 100"))
		}
		catalog, err := releaseCatalog(deps, info)
		if err != nil {
			return withExitCode(4, err)
		}
		start := (page - 1) * size
		if start > len(catalog.Releases) {
			start = len(catalog.Releases)
		}
		end := start + size
		if end > len(catalog.Releases) {
			end = len(catalog.Releases)
		}
		catalog.Releases = catalog.Releases[start:end]
		if format == "json" {
			return writeJSON(command, catalog)
		}
		fmt.Fprintln(command.OutOrStdout(), "VERSION  CHANNEL  PUBLISHED  CURRENT  LATEST  ASSET")
		for _, entry := range catalog.Releases {
			published := "-"
			if entry.PublishedAt != nil {
				if parsed, err := time.Parse(time.RFC3339, *entry.PublishedAt); err == nil {
					published = parsed.UTC().Format("2006-01-02")
				}
			}
			asset := "-"
			if entry.CompatibleAsset != nil {
				asset = entry.CompatibleAsset.Name
			}
			fmt.Fprintf(command.OutOrStdout(), "%s  %s  %s  %t  %t  %s\n", entry.Version, entry.Channel, published, entry.Current, entry.LatestStable, asset)
		}
		return nil
	}}
	command.Flags().IntVar(&page, "page", 1, "Select result page.")
	command.Flags().IntVar(&size, "size", 8, "Select page size.")
	command.Flags().StringVar(&format, "format", "text", "Select output format: text or json.")
	return command
}

func renderUpgrade(command *cobra.Command, envelope *updater.UpgradeEnvelope) {
	if envelope.Plan != nil {
		fmt.Fprintf(command.OutOrStdout(), "current: %s\ntarget: %s\nasset: %s\npath: %s\n", envelope.Plan.CurrentVersion, envelope.Plan.TargetVersion, envelope.Plan.AssetName, envelope.Plan.ExecutablePath)
	}
	fmt.Fprintf(command.OutOrStdout(), "outcome: %s\n", envelope.Outcome)
	if envelope.Error != nil {
		fmt.Fprintf(command.OutOrStdout(), "error: %s\n", envelope.Error.Message)
	}
	if envelope.Outcome == "failed" || envelope.Outcome == "rolled-back" {
		fmt.Fprintf(command.OutOrStdout(), "recovery: %s\nstop: %s\n", envelope.Recovery.InstallCommand, envelope.Recovery.StopProcessCommand)
	}
}

func newUninstallCommand(deps Dependencies) *cobra.Command {
	var dryRun bool
	var format string
	command := &cobra.Command{Use: "uninstall", Short: "Uninstall the native RunX executable.", Args: cobra.NoArgs, RunE: func(command *cobra.Command, _ []string) error {
		if err := validateFormat(format); err != nil {
			return err
		}
		path, err := deps.Executable()
		if err != nil {
			return withExitCode(5, err)
		}
		result := struct {
			Target      string `json:"target"`
			DryRun      bool   `json:"dryRun"`
			Uninstalled bool   `json:"uninstalled"`
		}{path, dryRun, false}
		if !dryRun {
			if err := os.Remove(path); err != nil {
				return withExitCode(5, fmt.Errorf("remove executable %s: %w", path, err))
			}
			result.Uninstalled = true
		}
		if format == "json" {
			return writeJSON(command, result)
		}
		if dryRun {
			fmt.Fprintf(command.OutOrStdout(), "would remove: %s\n", path)
		} else {
			fmt.Fprintf(command.OutOrStdout(), "removed: %s\n", path)
		}
		return nil
	}}
	command.Flags().BoolVar(&dryRun, "dry-run", false, "Print the target without deleting it.")
	command.Flags().StringVar(&format, "format", "text", "Select output format: text or json.")
	return command
}
