package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/CGuiho/runx/pkg/updater"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	dryRunFlag    bool
	targetVerFlag string
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Inspect or upgrade a native RunX executable.",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := updater.UpgradeOptions{
			DryRun:           dryRunFlag,
			CurrentVersion:   "0.8.0-dev",
			RequestedVersion: targetVerFlag,
		}

		if !jsonFmt {
			fmt.Println("------------------------------------------------------------")
			fmt.Println("  Upgrading the CLI")
			fmt.Println("------------------------------------------------------------")
			opts.OnPlan = func(plan updater.UpgradePlan) {
				fmt.Printf("  current : %s\n", plan.CurrentVersion)
				fmt.Printf("  target  : %s\n", plan.TargetVersion)
				fmt.Printf("  os      : %s\n", plan.OS)
				fmt.Printf("  arch    : %s\n", plan.Arch)
				if plan.AssetName != "" {
					fmt.Printf("  binary  : %s\n", plan.AssetName)
				}
				fmt.Printf("  path    : %s\n", plan.ExecutablePath)
				if plan.AssetURL != "" {
					fmt.Printf("  url     : %s\n", plan.AssetURL)
				}
				fmt.Println("------------------------------------------------------------")
			}
			opts.OnEvent = func(event updater.UpgradeEvent) {
				if event.Status == "started" && event.Phase != "plan" {
					switch event.Phase {
					case "download":
						fmt.Println("Downloading...")
					case "validate":
						fmt.Println("Validating...")
					case "replace":
						fmt.Println("Replacing...")
					case "verify":
						fmt.Println("Verifying...")
					case "cleanup":
						fmt.Println("Cleaning up...")
					}
				}
			}
		}

		envelope, err := updater.UpgradeSelf(opts)
		if err != nil {
			return err
		}

		if jsonFmt || viper.GetBool("json") {
			data, err := json.MarshalIndent(envelope, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))
			return nil
		}

		switch envelope.Outcome {
		case "upgraded":
			fmt.Printf("Upgrade complete: %s -> %s\n", envelope.Plan.CurrentVersion, envelope.Plan.TargetVersion)
		case "up-to-date":
			fmt.Printf("Already up to date: %s\n", envelope.Result.InstalledVersion)
		case "dry-run":
			fmt.Printf("Dry run complete: %s -> %s\n", envelope.Plan.CurrentVersion, envelope.Plan.TargetVersion)
		case "rolled-back":
			fmt.Printf("Upgrade failed during %s; restored RunX %s.\n", envelope.Error.Phase, envelope.Result.InstalledVersion)
		default:
			fmt.Printf("Upgrade failed during %s.\n", envelope.Error.Phase)
		}

		if envelope.Outcome == "failed" || envelope.Outcome == "rolled-back" {
			os.Exit(1)
		}

		return nil
	},
}

func init() {
	upgradeCmd.Flags().StringVar(&targetVerFlag, "version", "", "Select an exact release version.")
	upgradeCmd.Flags().StringVar(&upgradeArch, "arch", "", "Select target architecture.")
	upgradeCmd.Flags().StringVar(&upgradeVariant, "variant", "", "Select x64 binary variant.")
	upgradeCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "Plan without mutation.")

	RootCmd.AddCommand(upgradeCmd)
}
